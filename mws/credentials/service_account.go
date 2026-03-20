package credentials

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"go.mws.cloud/go-sdk/mws/iam"
	"go.mws.cloud/go-sdk/pkg/clock"
	"go.mws.cloud/go-sdk/service/iam/client"
)

// ServiceAccountProvider is a service account credentials provider.
type ServiceAccountProvider struct {
	*provider

	saKey       iam.ServiceAccountKey
	tokenIssuer *client.IssueServiceAccountTokenSugared
}

// NewServiceAccountProvider creates a new service account credentials provider.
func NewServiceAccountProvider(
	saKey iam.ServiceAccountKey,
	tokenIssuer client.IssueServiceAccountToken,
	opts ...ServiceAccountProviderOption,
) Provider {
	p := &ServiceAccountProvider{
		provider:    newProvider(),
		saKey:       saKey,
		tokenIssuer: client.NewIssueServiceAccountTokenSugared(tokenIssuer),
	}
	for _, opt := range opts {
		opt(p)
	}
	p.logger = p.logger.Named("service_account_provider")
	p.initCache()
	return p
}

// Provide returns the corresponding credentials for the service account
// specified in the key.
func (p *ServiceAccountProvider) Provide(ctx context.Context) (Credentials, error) {
	id := p.saKey.ServiceAccount.String()
	return p.do(id, func() (Credentials, error) {
		return p.provide(ctx, id)
	})
}

// InvalidateCredentials invalidates cached credentials.
func (p *ServiceAccountProvider) InvalidateCredentials(context.Context) error {
	id := p.saKey.ServiceAccount.String()
	return p.invalidateCredentials(id)
}

// Close closes the provider underlaying resources.
func (p *ServiceAccountProvider) Close(ctx context.Context) error {
	return p.closeCache(ctx)
}

func (p *ServiceAccountProvider) provide(ctx context.Context, id string) (Credentials, error) {
	signed, err := p.saKey.AuthorizedKey.JWS(p.clock.Now, id)
	if err != nil {
		return Credentials{}, err
	}

	req := client.IssueServiceAccountTokenV2Request{
		Authorization:  &signed,
		ServiceAccount: &id,
	}
	token, err := p.tokenIssuer.IssueServiceAccountTokenV2(ctx, req)
	if err != nil {
		return Credentials{}, fmt.Errorf("issue token: %w", err)
	}

	return Credentials{
		AccessToken: token.GetAccessToken(),
		ExpiresAt:   token.GetExpirationTsOr(p.clock.Now().Add(DefaultTokenTTL)),
	}, nil
}

// ServiceAccountProviderOption is a functional option for the service account
// credentials provider.
type ServiceAccountProviderOption func(*ServiceAccountProvider)

// WithServiceAccountProviderClock sets the clock for the service account
// credentials provider.
func WithServiceAccountProviderClock(clock clock.Clock) ServiceAccountProviderOption {
	return func(p *ServiceAccountProvider) {
		p.clock = clock
	}
}

// WithServiceAccountProviderLogger sets the logger for the service account
// credentials provider.
func WithServiceAccountProviderLogger(logger *zap.Logger) ServiceAccountProviderOption {
	return func(p *ServiceAccountProvider) {
		p.logger = logger
	}
}

// WithServiceAccountProviderCache sets the cache for the service account
// credentials provider.
func WithServiceAccountProviderCache(cache Cache) ServiceAccountProviderOption {
	return func(p *ServiceAccountProvider) {
		p.cache = cache
	}
}
