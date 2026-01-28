package mws

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"go.mws.cloud/util-toolset/pkg/os/env"
	"go.uber.org/zap"

	commonclient "go.mws.cloud/go-sdk/internal/client"
	"go.mws.cloud/go-sdk/mws/credentials"
	"go.mws.cloud/go-sdk/mws/endpoints"
	"go.mws.cloud/go-sdk/mws/iam"
	"go.mws.cloud/go-sdk/mws/retry"
	mwshttp "go.mws.cloud/go-sdk/pkg/http"
	iamclient "go.mws.cloud/go-sdk/service/iam/client/http"
)

// LoadSDKOption is a functional option that sets SDK load options.
type LoadSDKOption func(*loadSDKOptions)

// WithEnv sets the environment loader.
func WithEnv(env env.Env) LoadSDKOption {
	return func(o *loadSDKOptions) {
		o.env = env
	}
}

// WithDefaultProject sets a default project.
func WithDefaultProject(project string) LoadSDKOption {
	return func(c *loadSDKOptions) {
		c.defaultProject = project
	}
}

// WithDefaultZone sets a default zone.
func WithDefaultZone(zone string) LoadSDKOption {
	return func(c *loadSDKOptions) {
		c.defaultZone = zone
	}
}

// WithLogger sets the logger.
func WithLogger(logger *zap.Logger) LoadSDKOption {
	return func(o *loadSDKOptions) {
		o.logger = logger
	}
}

// WithConfig sets the SDK config.
func WithConfig(config Config) LoadSDKOption {
	return func(o *loadSDKOptions) {
		o.config = &config
	}
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(client HTTPClient) LoadSDKOption {
	return func(o *loadSDKOptions) {
		o.client = client
	}
}

// WithHTTPTransport sets the HTTP transport. Has no effect if [WithHTTPClient]
// is used.
func WithHTTPTransport(transport http.RoundTripper) LoadSDKOption {
	return func(o *loadSDKOptions) {
		o.transport = transport
	}
}

// WithTimeout sets the timeout for all client requests.
func WithTimeout(timeout time.Duration) LoadSDKOption {
	return func(o *loadSDKOptions) {
		o.timeout = timeout
	}
}

// WithRetryer sets the retryer.
func WithRetryer(retryer retry.Retryer) LoadSDKOption {
	return func(o *loadSDKOptions) {
		o.retryer = retryer
	}
}

// WithServiceEndpointResolver sets the service endpoint resolver.
func WithServiceEndpointResolver(resolver endpoints.ServiceEndpointResolver) LoadSDKOption {
	return func(o *loadSDKOptions) {
		o.serviceEndpointResolver = resolver
	}
}

// WithCredentials sets the credentials provider.
func WithCredentials(provider credentials.Provider) LoadSDKOption {
	return func(c *loadSDKOptions) {
		c.credentials = provider
	}
}

type loadSDKOptions struct {
	env                     env.Env
	config                  *Config
	defaultProject          string
	defaultZone             string
	timeout                 time.Duration
	logger                  *zap.Logger
	client                  HTTPClient
	transport               http.RoundTripper
	retryer                 retry.Retryer
	serviceEndpointResolver endpoints.ServiceEndpointResolver
	credentials             credentials.Provider
}

func newLoadSDKOptions() *loadSDKOptions {
	return &loadSDKOptions{
		env:       env.RealEnv{},
		transport: mwshttp.DefaultTransport(),
	}
}

func (o *loadSDKOptions) build(ctx context.Context) (*SDK, error) {
	sdk := &SDK{}
	o.setOpts(sdk)
	if err := o.setDefaults(ctx, sdk); err != nil {
		return nil, err
	}
	return sdk, nil
}

func (o *loadSDKOptions) setOpts(sdk *SDK) {
	sdk.defaultProject = o.defaultProject
	sdk.defaultZone = o.defaultZone
	if o.logger != nil {
		sdk.logger = o.logger
	}
	if o.client != nil {
		sdk.client = o.client
	}
	if o.retryer != nil {
		sdk.retryer = o.retryer
	}
	if o.serviceEndpointResolver != nil {
		sdk.serviceEndpointResolver = o.serviceEndpointResolver
	}
	if o.credentials != nil {
		sdk.credentials = o.credentials
	}
}

func (o *loadSDKOptions) setDefaults(ctx context.Context, sdk *SDK) (err error) {
	if o.config == nil {
		if o.config, err = LoadConfig(LoadConfigWithEnv(o.env)); err != nil {
			return fmt.Errorf("load config: %w", err)
		}
	}
	if sdk.defaultProject == "" {
		sdk.defaultProject = o.config.Project
	}
	if sdk.defaultZone == "" {
		sdk.defaultZone = o.config.Zone
	}
	if sdk.logger == nil {
		sdk.logger, err = o.buildLogger(o.config.LogLevel)
		if err != nil {
			return fmt.Errorf("build logger: %w", err)
		}
	}
	if sdk.client == nil {
		sdk.client = o.buildClient()
	}

	sdk.client, sdk.clientCancel = newCancelableClient(sdk.client)

	if sdk.retryer == nil {
		sdk.retryer = retry.NewStandardRetry()
	}
	if sdk.serviceEndpointResolver == nil {
		sdk.serviceEndpointResolver = o.buildServiceEndpointResolver(sdk.logger, sdk.client)
	}
	if sdk.credentials == nil {
		if sdk.credentials, err = o.buildCredentials(ctx, sdk); err != nil {
			return fmt.Errorf("build credentials: %w", err)
		}
	}
	return nil
}

func (o *loadSDKOptions) buildLogger(level string) (*zap.Logger, error) {
	if level == "" {
		return zap.NewNop(), nil
	}

	lvl, err := zap.ParseAtomicLevel(o.config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Named("sdk"), nil
}

func (o *loadSDKOptions) buildClient() *http.Client {
	return &http.Client{
		Timeout: o.getTimeout(),
		Transport: mwshttp.Chain(
			mwshttp.RequestIDInjector,
			mwshttp.IdempotencyTokenInjector,
			mwshttp.RetryAttemptInjector,
			mwshttp.UserAgentInjector(o.userAgent()),
		)(o.transport),
	}
}

func (o *loadSDKOptions) getTimeout() time.Duration {
	return cmp.Or(o.timeout, o.config.Timeout, DefaultTimeout)
}

func (o *loadSDKOptions) userAgent() string {
	version := "unknown"

	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		version = info.Main.Version
	}

	return "mws/go-sdk/" + version
}

func (o *loadSDKOptions) buildCredentials(ctx context.Context, sdk *SDK) (credentials.Provider, error) {
	if o.config.Token != "" {
		return credentials.StaticProvider(credentials.Credentials{
			AccessToken: o.config.Token,
		}), nil
	}

	if o.config.ServiceAccountAuthorizedKeyPath != "" {
		return o.buildServiceAccountAuthorizedKeyCredentials(ctx, sdk)
	}

	return credentials.AnonymousProvider(), nil
}

func (o *loadSDKOptions) buildServiceAccountAuthorizedKeyCredentials(ctx context.Context, sdk *SDK) (credentials.Provider, error) {
	var key iam.ServiceAccountKey

	data, err := os.ReadFile(o.config.ServiceAccountAuthorizedKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read service account key authorized file: %w", err)
	}

	if err = json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("unmarshal service account authorized key: %w", err)
	}

	endpoint, err := sdk.serviceEndpointResolver.Resolve(ctx, "iam")
	if err != nil {
		return nil, fmt.Errorf("resolve iam service endpoint: %w", err)
	}

	config, err := commonclient.NewClientConfig(string(endpoint), commonclient.WithHTTPClient(sdk.client))
	if err != nil {
		return nil, fmt.Errorf("create iam client config: %w", err)
	}

	tokenIssuer := iamclient.NewIssueServiceAccountToken(*config)
	return credentials.NewServiceAccountProvider(key, tokenIssuer,
		credentials.WithServiceAccountProviderLogger(sdk.logger),
	), nil
}

func (o *loadSDKOptions) buildServiceEndpointResolver(logger *zap.Logger, client HTTPClient) endpoints.ServiceEndpointResolver {
	endpoint := endpoints.Endpoint(o.config.BaseEndpoint)
	discoveryClient := endpoints.NewHTTPDiscoveryClient(logger, client, endpoint)
	return endpoints.NewDiscoveryServiceEndpointResolver(discoveryClient)
}

type httpClientFunc func(*http.Request) (*http.Response, error)

func (fn httpClientFunc) Do(r *http.Request) (*http.Response, error) { return fn(r) }

func newCancelableClient(client HTTPClient) (HTTPClient, func()) {
	if client == nil {
		panic("HTTPClient is nil")
	}

	clientCtx, clientCancel := context.WithCancel(context.Background())

	result := httpClientFunc(func(request *http.Request) (*http.Response, error) {
		if clientCtx.Err() != nil {
			return nil, ErrSDKClosed
		}

		ctx, cancel := context.WithCancelCause(request.Context())

		context.AfterFunc(clientCtx, func() { cancel(ErrSDKClosed) })

		result, err := client.Do(request.WithContext(ctx))
		switch {
		case err == nil:
			return result, nil
		case errors.Is(err, context.Canceled):
			return nil, context.Cause(ctx)
		default:
			return nil, err
		}
	})

	closer, ok := client.(closeIdler)
	if !ok {
		closer = noopCloser{}
	}

	return &struct {
		HTTPClient
		closeIdler
	}{result, closer}, clientCancel
}

type closeIdler interface{ CloseIdleConnections() }

type noopCloser struct{}

func (noopCloser) CloseIdleConnections() {}
