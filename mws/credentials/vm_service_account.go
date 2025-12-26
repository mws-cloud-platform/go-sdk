package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"go.mws.cloud/go-sdk/pkg/clock"
)

const (
	tokenKey = "instance/service-accounts/default/token"
	nameKey  = "instance/name"
)

// MetadataProvider represents VM instance metadata provider.
//
// It's compatible with Google Compute Engine (GCE) metadata service provider
// implementation: https://pkg.go.dev/cloud.google.com/go/compute/metadata.
type MetadataProvider interface {
	// GetWithContext returns a value from the metadata service.
	GetWithContext(ctx context.Context, key string) (string, error)
}

// VMServiceAccountProvider is a credentials provider that issues credentials
// for the service account connected to the running VM.
type VMServiceAccountProvider struct {
	*provider

	metadataProvider MetadataProvider
}

// NewVMServiceAccountProvider creates a new VM service account credentials
// provider.
func NewVMServiceAccountProvider(
	metadataProvider MetadataProvider,
	opts ...VMServiceAccountProviderOption,
) *VMServiceAccountProvider {
	p := &VMServiceAccountProvider{
		provider:         newProvider(),
		metadataProvider: metadataProvider,
	}
	for _, opt := range opts {
		opt(p)
	}
	p.logger = p.logger.Named("vm_service_account_provider")
	p.initCache()
	return p
}

// Provide returns credentials for the service account connected to the VM.
func (p *VMServiceAccountProvider) Provide(ctx context.Context) (Credentials, error) {
	return p.do(tokenKey, func() (Credentials, error) {
		return p.provide(ctx)
	})
}

// InvalidateCredentials invalidates cached credentials.
func (p *VMServiceAccountProvider) InvalidateCredentials(context.Context) error {
	return p.invalidateCredentials(tokenKey)
}

// Close closes the provider underlaying resources.
func (p *VMServiceAccountProvider) Close(ctx context.Context) error {
	return p.closeCache(ctx)
}

func (p *VMServiceAccountProvider) provide(ctx context.Context) (Credentials, error) {
	tokenJSON, err := p.metadataProvider.GetWithContext(ctx, tokenKey)
	if err != nil {
		return Credentials{}, err
	}

	var token struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err = json.Unmarshal([]byte(tokenJSON), &token); err != nil {
		return Credentials{}, fmt.Errorf("unmarshal token: %w", err)
	}

	p.logger.Info("token received", zap.Int64("access_token_expires_in", token.ExpiresIn))

	return Credentials{
		AccessToken: token.AccessToken,
		ExpiresAt:   p.clock.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	}, nil
}

// VMServiceAccountProviderOption is a functional option for the VM service
// account credentials provider.
type VMServiceAccountProviderOption func(*VMServiceAccountProvider)

// WithVMServiceAccountProviderClock sets the clock for the VM service account
// credentials provider.
func WithVMServiceAccountProviderClock(clock clock.Clock) VMServiceAccountProviderOption {
	return func(p *VMServiceAccountProvider) {
		p.clock = clock
	}
}

// WithVMServiceAccountProviderLogger sets the logger for the VM service account
// credentials provider.
func WithVMServiceAccountProviderLogger(logger *zap.Logger) VMServiceAccountProviderOption {
	return func(p *VMServiceAccountProvider) {
		p.logger = logger
	}
}

// WithVMServiceAccountProviderTokenExpirationDelta sets the token expiration
// delta for the VM service account credentials provider.
func WithVMServiceAccountProviderTokenExpirationDelta(delta time.Duration) VMServiceAccountProviderOption {
	return func(p *VMServiceAccountProvider) {
		p.tokenExpirationDelta = delta
	}
}

// WithVMServiceAccountProviderCache sets the cache for the VM service account
// credentials provider.
func WithVMServiceAccountProviderCache(cache Cache) VMServiceAccountProviderOption {
	return func(p *VMServiceAccountProvider) {
		p.cache = cache
	}
}

// OnComputeVM returns true if this code is running on a Compute VM.
// Note: true from this function doesn't guarantee that all the metadata is defined.
func OnComputeVM(ctx context.Context, metadataProvider MetadataProvider) bool {
	result, err := metadataProvider.GetWithContext(ctx, nameKey)
	return err == nil && result != ""
}
