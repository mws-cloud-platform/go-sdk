package credentials

import (
	"context"
	"sync"

	"go.mws.cloud/go-sdk/pkg/clock"
)

type DumbCache struct {
	mu   sync.RWMutex
	data map[string]Credentials

	clock clock.Clock
}

type DumbCacheOption func(*DumbCache)

func WithClockCache(clock clock.Clock) DumbCacheOption {
	return func(c *DumbCache) {
		c.clock = clock
	}
}

func NewDumbCache(opt ...DumbCacheOption) *DumbCache {
	cm := &DumbCache{
		data:  make(map[string]Credentials),
		clock: clock.NewReal(),
	}
	for _, o := range opt {
		o(cm)
	}
	return cm
}

func (c *DumbCache) Load(key string) (Credentials, error) {
	c.mu.RLock()
	creds, ok := c.data[key]
	c.mu.RUnlock()
	if !ok {
		return Credentials{}, ErrEntryNotFound
	}
	now := c.clock.Now()
	if creds.ExpiresAt.Before(now) {
		c.mu.Lock()
		defer c.mu.Unlock()
		creds, ok = c.data[key]
		if !ok {
			return Credentials{}, ErrEntryNotFound
		}
		if creds.ExpiresAt.Before(now) {
			delete(c.data, key)
			creds = Credentials{}
			return Credentials{}, ErrEntryNotFound
		}
		return creds, nil
	}

	return creds, nil
}

func (c *DumbCache) Store(key string, creds Credentials) error {
	if creds.ExpiresAt.Before(c.clock.Now()) {
		return ErrEntryRejected
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = creds
	return nil
}

func (c *DumbCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *DumbCache) Close(context.Context) error {
	return nil
}

func (c *DumbCache) IsEmpty() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data) == 0
}
