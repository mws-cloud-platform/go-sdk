package client

import (
	"context"
	"net/http"
)

type Client struct {
	httpClient  HTTPClient
	interceptor Interceptor
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
	return c.httpClient.Do(r)
}

func (c *Client) Intercept(ctx context.Context, request any, response APIResp, invoker Invoker) error {
	if c.interceptor != nil {
		return c.interceptor(ctx, request, response, invoker)
	}

	return invoker(ctx, request, response)
}

type Invoker func(ctx context.Context, request any, response APIResp) error

type Interceptor func(ctx context.Context, request any, response APIResp, invoker Invoker) error

type APIResp interface {
	GetErr() error
	GetCode() int
}
