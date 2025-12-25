package credentials

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"go.mws.cloud/go-sdk/pkg/clock"
)

type provider struct {
	clock                clock.Clock
	logger               *zap.Logger
	singleflight         singleflight.Group
	cache                Cache
	tokenExpirationDelta time.Duration
	closeCache           func(context.Context) error
}

func newProvider() *provider {
	return &provider{
		clock:                clock.NewReal(),
		logger:               zap.NewNop(),
		tokenExpirationDelta: DefaultTokenExpirationDelta,
		closeCache:           func(context.Context) error { return nil },
	}
}

func (p *provider) initCache() {
	if p.cache == nil {
		cache := newTTLCache()
		p.cache = cache
		p.closeCache = cache.Close
	}
}

func (p *provider) do(key string, fn func() (Credentials, error)) (Credentials, error) {
	cached, cacheErr := p.cache.Load(key)
	switch {
	case cacheErr == nil && cached.ExpiresAt.After(p.clock.Now().Add(p.tokenExpirationDelta)):
		p.logger.Debug("using cached credentials", zap.Duration("token_expiration_delta", p.tokenExpirationDelta))
		return cached, nil
	case errors.Is(cacheErr, ErrEntryNotFound):
	case cacheErr != nil:
		return Credentials{}, cacheErr
	}

	res, err, _ := p.singleflight.Do(key, func() (any, error) { return fn() })
	if err != nil {
		// If there is a credentials in the cache, we will return it to the user in case of error in the passed function.
		if cacheErr == nil && cached.ExpiresAt.After(p.clock.Now()) {
			p.logger.Warn("using expiring credentials due to token issue error", zap.Error(err))
			return cached, nil
		}
		return Credentials{}, fmt.Errorf("issue token: %w", err)
	}
	creds := res.(Credentials)

	return creds, p.cache.Store(key, creds)
}

func (p *provider) invalidateCredentials(key string) error {
	return p.cache.Delete(key)
}
