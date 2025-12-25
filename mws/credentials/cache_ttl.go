package credentials

import (
	"context"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

const (
	cacheDefaultBufferItems = int64(64)
	cacheDefaultItemSize    = 1024
	cacheDefaultCapacity    = 1e4
	cacheUpdatePeriodInSec  = 1
	cacheEstimatedOverhead  = 32
)

// TTLCacheOption is a functional option for configuring a [TTLCache].
type TTLCacheOption func(*TTLCache)

// WithTTLCacheSync sets the synchronous/asynchronous operation mode:
//   - sync - element is available immediately after addition; lower performance
//   - async - element may not be immediately available after addition; higher performance
func WithTTLCacheSync(sync bool) TTLCacheOption {
	return func(c *TTLCache) {
		c.sync = sync
	}
}

// WithTTLCacheCapacity sets the maximum cache capacity, default - 1e4.
func WithTTLCacheCapacity(capacity int64) TTLCacheOption {
	return func(c *TTLCache) {
		c.capacity = capacity
	}
}

// WithTTLCacheItemSize sets the size of each item in the cache, default - 1024.
func WithTTLCacheItemSize(itemSize int64) TTLCacheOption {
	return func(c *TTLCache) {
		c.itemSize = itemSize
	}
}

// WithTTLCacheCleanPeriod sets the cache clean period. Rounded down to the
// nearest second. Cannot be less than 1 second. Default - 1 second.
func WithTTLCacheCleanPeriod(every time.Duration) TTLCacheOption {
	return func(c *TTLCache) {
		if every < time.Second {
			every = time.Second
		}
		c.ttlCleanPeriodInSec = int64(every / time.Second)
	}
}

// TTLCache is a credentials cache with TTL support.
type TTLCache struct {
	cache               *ristretto.Cache[string, Credentials]
	sync                bool
	itemSize            int64
	capacity            int64
	bufferItems         int64
	ttlCleanPeriodInSec int64
}

// NewTTLCache creates a new [TTLCache] instance with the specified options.
func NewTTLCache(opts ...TTLCacheOption) (c *TTLCache, err error) {
	c = &TTLCache{
		sync:                false,
		itemSize:            cacheDefaultItemSize,
		capacity:            cacheDefaultCapacity,
		bufferItems:         cacheDefaultBufferItems,
		ttlCleanPeriodInSec: cacheUpdatePeriodInSec,
	}
	for _, o := range opts {
		o(c)
	}
	numCounters, maxCost := calculateOptimalConfig(c.capacity, c.itemSize)
	c.cache, err = ristretto.NewCache(&ristretto.Config[string, Credentials]{
		NumCounters:            numCounters,
		MaxCost:                maxCost,
		BufferItems:            c.bufferItems,
		TtlTickerDurationInSec: c.ttlCleanPeriodInSec,
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Load loads credentials from the cache by key.
func (c *TTLCache) Load(key string) (Credentials, error) {
	cred, ok := c.cache.Get(key)
	if !ok {
		return cred, ErrEntryNotFound
	}
	return cred, nil
}

// Store stores credentials in the cache with a TTL based on the token
// expiration time.
func (c *TTLCache) Store(key string, credentials Credentials) error {
	ok := c.cache.SetWithTTL(key, credentials, c.itemSize, time.Until(credentials.ExpiresAt))
	if !ok {
		return ErrEntryRejected
	}
	if c.sync {
		c.cache.Wait()
	}
	return nil
}

// Delete deletes credentials by key from the cache if it exists.
func (c *TTLCache) Delete(key string) error {
	c.cache.Del(key)
	if c.sync {
		c.cache.Wait()
	}
	return nil
}

// Close stops all goroutines and closes all channels.
func (c *TTLCache) Close(context.Context) error {
	c.cache.Close()
	return nil
}

// Based on the library documentation: https://pkg.go.dev/github.com/dgraph-io/ristretto/v2#Config.
func calculateOptimalConfig(expectedMaxItems, avgValueSize int64) (int64, int64) {
	numCounters := expectedMaxItems * 10

	maxCost := expectedMaxItems * (avgValueSize + cacheEstimatedOverhead)

	return numCounters, maxCost
}

func newTTLCache(opts ...TTLCacheOption) *TTLCache {
	c, err := NewTTLCache(opts...)
	if err != nil {
		panic(err)
	}
	return c
}
