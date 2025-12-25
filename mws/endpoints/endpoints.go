// Package endpoints provides resolvers for service API endpoints.
package endpoints

import (
	"context"
	"fmt"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

const (
	// ErrEndpointNotFound reports whether no endpoint was found for a service.
	ErrEndpointNotFound = consterr.Error("endpoint not found")
)

// Endpoint represents the API endpoint.
type Endpoint string

// ServiceName represents the service name. Each SDK service client should have
// its own service name constant.
type ServiceName string

// ServiceEndpointResolver represents a resolver for service API endpoint.
type ServiceEndpointResolver interface {
	Resolve(context.Context, ServiceName) (Endpoint, error)
}

// Chain chains resolvers together. The returned resolver subsequently tries to
// get the service API endpoint from the given resolvers. If resolver succeeds
// the obtained endpoint is returned. If resolver fails, the next resolver is
// tried.
func Chain(resolvers ...ServiceEndpointResolver) ServiceEndpointResolver {
	return chain{resolvers: resolvers}
}

type chain struct {
	resolvers []ServiceEndpointResolver
}

func (c chain) Resolve(ctx context.Context, service ServiceName) (Endpoint, error) {
	for _, r := range c.resolvers {
		endpoint, err := r.Resolve(ctx, service)
		if err != nil {
			continue
		}
		return endpoint, nil
	}
	return "", errServiceEndpointNotFound(service)
}

func errServiceEndpointNotFound(service ServiceName) error {
	return fmt.Errorf("service %q: %w", service, ErrEndpointNotFound)
}
