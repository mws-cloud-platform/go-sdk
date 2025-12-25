package client

import (
	"context"
)

// DefaultsInjector creates a client interceptor that injects default values
// into requests.
func DefaultsInjector(project, zone string) Interceptor {
	return func(ctx context.Context, request any, response APIResp, invoker Invoker) error {
		if project != "" {
			if req, ok := request.(requestWithProject); ok && req.GetProject() == "" {
				req.SetProject(project)
			}
		}
		if zone != "" {
			if req, ok := request.(requestWithZone); ok && req.GetZone() == "" {
				req.SetZone(zone)
			}
		}
		return invoker(ctx, request, response)
	}
}

type requestWithProject interface {
	SetProject(string)
	GetProject() string
}

type requestWithZone interface {
	SetZone(string)
	GetZone() string
}
