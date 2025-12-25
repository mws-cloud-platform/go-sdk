package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.mws.cloud/util-toolset/pkg/net/http/client"
	"go.uber.org/zap"
)

const defaultDiscoveryPath = "/endpoint"

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// DiscoveryServiceEndpointResolver is a service API endpoint resolver which gets
// service API endpoints from the discovery endpoint.
type DiscoveryServiceEndpointResolver struct {
	client DiscoveryClient
}

// NewDiscoveryServiceEndpointResolver creates a new discovery service API endpoint resolver.
func NewDiscoveryServiceEndpointResolver(client DiscoveryClient) ServiceEndpointResolver {
	return DiscoveryServiceEndpointResolver{
		client: client,
	}
}

// Resolve returns the service API endpoint from the discovery response.
// Service may be in any format, but if it has slashes, e.g. "iam/authservice",
// then only the part before the first slash will be used, in this case it will be "iam".
// If service is not found in the mapping, an error is returned.
func (s DiscoveryServiceEndpointResolver) Resolve(ctx context.Context, service ServiceName) (Endpoint, error) {
	endpoints, err := s.client.Endpoints(ctx)
	if err != nil {
		return "", fmt.Errorf("get endpoints list: %w", err)
	}

	key, _, _ := strings.Cut(string(service), "/")
	endpoint, ok := endpoints[key]
	if !ok {
		return "", errServiceEndpointNotFound(service)
	}
	return Endpoint(endpoint.Address), nil
}

// DiscoveryClient represents a discovery client that can fetch endpoints
// information from the service discovery.
type DiscoveryClient interface {
	Endpoints(context.Context) (DiscoveryEndpoints, error)
}

// HTTPDiscoveryClient is a discovery client implementation that fetches
// endpoints information from the service discovery over HTTP.
type HTTPDiscoveryClient struct {
	logger   *zap.Logger
	client   HTTPClient
	endpoint string
}

// NewHTTPDiscoveryClient creates a new discovery client.
func NewHTTPDiscoveryClient(logger *zap.Logger, client HTTPClient, endpoint Endpoint) HTTPDiscoveryClient {
	return HTTPDiscoveryClient{
		logger:   logger,
		client:   client,
		endpoint: string(endpoint),
	}
}

// Endpoints fetches endpoints information from the service discovery.
func (c HTTPDiscoveryClient) Endpoints(ctx context.Context) (DiscoveryEndpoints, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimSuffix(c.endpoint, "/")+defaultDiscoveryPath, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			c.logger.Warn("close response body", zap.Error(closeErr))
		}
	}()

	var endpoints []discoveryEndpointJSON
	if err = client.ReadJSON(response.Body, &endpoints); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	m := make(DiscoveryEndpoints, len(endpoints))
	for _, e := range endpoints {
		m[e.ID] = DiscoveryEndpoint{Address: e.Address}
	}
	return m, nil
}

// DiscoveryEndpoint contains service discovery endpoint information.
type DiscoveryEndpoint struct {
	Address string
}

// DiscoveryEndpoints contains service discovery endpoints information grouped
// by service names.
type DiscoveryEndpoints map[string]DiscoveryEndpoint

type discoveryEndpointJSON struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}
