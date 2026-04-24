package credentials_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"go.mws.cloud/go-sdk/mws/credentials"
)

const cacheItemCount = int(1e3)

type cacheCloser interface {
	credentials.Cache
	Close(ctx context.Context) error
}

func TestTTLCache_StoreAndLoad(t *testing.T) {
	defer goleak.VerifyNone(t)
	ttlcache, err := credentials.NewTTLCache(credentials.WithTTLCacheSync(true))
	require.NoError(t, err)
	testCases := []struct {
		name  string
		cache cacheCloser
	}{
		{name: "Real cache", cache: ttlcache},
		{name: "Mock cache", cache: credentials.NewDumbCache()},
	}
	for _, data := range testCases {
		t.Run(data.name, func(t *testing.T) {
			cache := data.cache
			defer func() {
				require.NoError(t, cache.Close(t.Context()))
			}()
			creds := generateCredentials(time.Hour)
			err = fillCache(cache, creds)
			require.NoError(t, err)
			for key, cred := range creds {
				loadedCreds, err := cache.Load(key)
				require.NoError(t, err)
				require.Equal(t, cred, loadedCreds)
			}
		})
	}
}

func TestTTLCache_LoadNotFound(t *testing.T) {
	defer goleak.VerifyNone(t)
	ttlcache, err := credentials.NewTTLCache(credentials.WithTTLCacheSync(true))
	require.NoError(t, err)
	testCases := []struct {
		name  string
		cache cacheCloser
	}{
		{name: "Real cache", cache: ttlcache},
		{name: "Mock cache", cache: credentials.NewDumbCache()},
	}
	for _, data := range testCases {
		t.Run(data.name, func(t *testing.T) {
			cache := data.cache
			defer func() {
				require.NoError(t, cache.Close(t.Context()))
			}()
			_, err = cache.Load("not_found_key")
			require.ErrorIs(t, err, credentials.ErrEntryNotFound)
		})
	}
}

func TestTTLCache_Remove(t *testing.T) {
	defer goleak.VerifyNone(t)
	ttlcache, err := credentials.NewTTLCache(credentials.WithTTLCacheSync(true))
	require.NoError(t, err)
	testCases := []struct {
		name  string
		cache cacheCloser
	}{
		{name: "Real cache", cache: ttlcache},
		{name: "Mock cache", cache: credentials.NewDumbCache()},
	}
	for _, data := range testCases {
		t.Run(data.name, func(t *testing.T) {
			cache := data.cache
			defer func() {
				require.NoError(t, cache.Close(t.Context()))
			}()
			creds := generateCredentials(time.Hour)
			err = fillCache(cache, creds)
			require.NoError(t, err)
			i := 0
			for key, cred := range creds {
				if i < cacheItemCount/2 { // Удаляем только первую половину
					err = cache.Delete(key)
					require.NoError(t, err)
					_, err := cache.Load(key)
					require.ErrorIs(t, err, credentials.ErrEntryNotFound)
				} else {
					loadedCreds, err := cache.Load(key)
					require.NoError(t, err)
					require.Equal(t, cred, loadedCreds)
				}
				i++
			}
		})
	}
}

func TestTTLCache_ExpireTTL(t *testing.T) {
	testCases := []struct {
		name  string
		cache func(t *testing.T) cacheCloser
	}{
		{name: "Real cache", cache: func(t *testing.T) cacheCloser {
			ttlcache, err := credentials.NewTTLCache(
				credentials.WithTTLCacheCleanPeriod(time.Minute),
				credentials.WithTTLCacheSync(true),
			)
			require.NoError(t, err)
			return ttlcache
		}},
		{name: "Mock cache", cache: func(*testing.T) cacheCloser { return credentials.NewDumbCache() }},
	}
	for _, data := range testCases {
		t.Run(data.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				cache := data.cache(t)
				defer func() {
					require.NoError(t, cache.Close(t.Context()))
				}()
				creds := generateCredentials(4 * time.Hour)
				credsTTL := generateCredentials(time.Hour)
				err := fillCache(cache, creds)
				require.NoError(t, err)
				err = fillCache(cache, credsTTL)
				require.NoError(t, err)
				synctest.Wait()
				t.Log("initial load complete, sleeping for 1 hour to test TTL expiration", time.Now())
				time.Sleep(time.Hour)
				synctest.Wait()
				for key, cred := range credsTTL { // Проверяем креды с TTL, должны не истечь
					loadedCreds, err := cache.Load(key)
					require.NoError(t, err)
					require.Equal(t, cred, loadedCreds)
				}
				t.Log("1 hour sleep complete, sleeping for another minute to ensure entries expulsion", time.Now())
				time.Sleep(time.Minute) // Ждем истечения TTL для части кредов
				t.Log("additional sleep complete, verifying cache entries expulsion", time.Now())
				for key, cred := range creds { // Проверяем обычные креды
					loadedCreds, err := cache.Load(key)
					require.NoError(t, err)
					require.Equal(t, cred, loadedCreds)
				}
				for key := range credsTTL { // Проверяем гарантировано истекшие креды
					_, err := cache.Load(key)
					require.ErrorIs(t, err, credentials.ErrEntryNotFound)
				}
			})
		})
	}
}

func generateCredentials(ttl time.Duration) map[string]credentials.Credentials {
	res := make(map[string]credentials.Credentials)
	expires := time.Now().Add(ttl)
	for range cacheItemCount {
		strUUID := uuid.NewString()
		res[strUUID] = credentials.Credentials{AccessToken: strUUID, ExpiresAt: expires}
	}
	return res
}

func fillCache(cache credentials.Cache, creds map[string]credentials.Credentials) error {
	for k, v := range creds {
		if err := cache.Store(k, v); err != nil {
			return err
		}
	}
	return nil
}
