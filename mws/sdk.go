// Package mws provides the SDK and related types and utilities.
package mws

import (
	"context"
	"net/http"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.uber.org/zap"

	commonclient "go.mws.cloud/go-sdk/internal/client"
	authinterceptor "go.mws.cloud/go-sdk/internal/client/interceptors/auth"
	loginterceptor "go.mws.cloud/go-sdk/internal/client/interceptors/log"
	recoveryinterceptor "go.mws.cloud/go-sdk/internal/client/interceptors/recovery"
	retryinterceptor "go.mws.cloud/go-sdk/internal/client/interceptors/retry"
	"go.mws.cloud/go-sdk/mws/credentials"
	"go.mws.cloud/go-sdk/mws/endpoints"
	"go.mws.cloud/go-sdk/mws/retry"
)

// ErrSDKClosed is returned when the [SDK] is closed.
const ErrSDKClosed = consterr.Error("sdk is closed")

// HTTPClient represents an HTTP client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// SDK provides configuration for service clients.
//
// Use [Load] to create SDK instance.
type SDK struct {
	logger                  *zap.Logger
	defaultZone             string
	defaultProject          string
	client                  HTTPClient
	clientCancel            func()
	retryer                 retry.Retryer
	serviceEndpointResolver endpoints.ServiceEndpointResolver
	credentials             credentials.Provider
}

// Load loads SDK from the specified configuration.
//
//	sdk, err := mws.Load(context.TODO())
//	if err != nil {
//		log.Fatal("load sdk:", err)
//	}
//
// By default, SDK loads configuration from the environment variables and
// sensible defaults. Check the [Config] for more information about SDK
// configuration parameters. Provide [LoadSDKOption] options to override default
// behavior.
func Load(ctx context.Context, opts ...LoadSDKOption) (*SDK, error) {
	options := newLoadSDKOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options.build(ctx)
}

// DefaultProject returns a default project used by the SDK.
func (s *SDK) DefaultProject() string {
	return s.defaultProject
}

// DefaultZone returns a default zone used by the SDK.
func (s *SDK) DefaultZone() string {
	return s.defaultZone
}

// ServiceEndpointResolver returns service API endpoint resolver.
func (s *SDK) ServiceEndpointResolver() endpoints.ServiceEndpointResolver {
	return s.serviceEndpointResolver
}

// CredentialsProvider returns credentials provider.
func (s *SDK) CredentialsProvider() credentials.Provider {
	return s.credentials
}

// ClientOptions returns service client options based on the SDK configuration.
func (s *SDK) ClientOptions() []commonclient.Option {
	return []commonclient.Option{
		commonclient.WithHTTPClient(s.client),
		commonclient.WithChainInterceptor(
			recoveryinterceptor.Recovery,
			commonclient.RequestIDInjector,
			commonclient.IdempotencyKeyInjector,
			commonclient.DefaultsInjector(s.defaultProject, s.defaultZone),
			authinterceptor.New(s.credentials),
			loginterceptor.New(s.logger.Named("client")),
			retryinterceptor.New(retryinterceptor.WithRetryer(s.retryer)),
			commonclient.ErrorWrapper,
		),
	}
}

// Close closes the [SDK], canceling all in-flight requests. It will attempt to
// also close the idle connections of the underlying HTTP client if it supports
// it. It will also stop background goroutines and release resources.
//
// If it's called concurrently with in-flight requests in other goroutines,
// those requests will be canceled, but cached [net] connections themselves may
// prevent [http.Transport] from being garbage collected, resulting in a memory
// leak, since there is no guarantee that CloseIdleConnections will be called
// after all connections had returned to the pool. You can avert this by
// ensuring that all requests are finished before calling this method. Or you
// can pass your own [HTTPClient] using [WithHTTPClient] and manage
// [http.Transport] lifecycle yourself.
func (s *SDK) Close(ctx context.Context) error {
	s.clientCancel()
	if c, ok := s.client.(closeIdler); ok {
		c.CloseIdleConnections()
	}
	c, ok := s.credentials.(credentials.Closer)
	if !ok {
		return nil
	}
	return c.Close(ctx)
}
