// Package credentials provides utilities for credentials retrieval and
// management.
package credentials

import (
	"context"
	"time"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

const (
	// ErrUnauthorized reports whether the user is not authorized.
	ErrUnauthorized = consterr.Error("unauthorized")

	// DefaultTokenTTL is the default token TTL.
	DefaultTokenTTL = 60 * time.Minute

	// DefaultTokenExpirationDelta is the default token expiration delta.
	DefaultTokenExpirationDelta = 5 * time.Minute
)

// Provider represents a provider of credentials for service clients.
type Provider interface {
	Provide(context.Context) (Credentials, error)
}

// Closer represents a closer of credentials provider. Provider can optionally
// implement this interface.
type Closer interface {
	Close(context.Context) error
}

// ProviderFunc is an adapter to allow the use of ordinary functions as
// credentials provider.
type ProviderFunc func(context.Context) (Credentials, error)

func (f ProviderFunc) Provide(ctx context.Context) (Credentials, error) {
	return f(ctx)
}

// Credentials contains credentials for service clients.
type Credentials struct {
	AccessToken string
	ExpiresAt   time.Time
}

// StaticProvider is a static credentials provider.
func StaticProvider(creds Credentials) Provider {
	return ProviderFunc(func(context.Context) (Credentials, error) {
		return creds, nil
	})
}

// AnonymousProvider is an anonymous credentials provider. Returns an empty credentials.
func AnonymousProvider() Provider {
	return ProviderFunc(func(context.Context) (Credentials, error) {
		return Credentials{}, nil
	})
}
