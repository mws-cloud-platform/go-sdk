package client

import (
	"context"
	"net/http"
)

type MWSClient interface {
	HTTPClient
	Intercept(ctx context.Context, request any, response APIResp, invoker Invoker) error
}

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}
