package client

import (
	"context"
	"net/http"
	"net/url"
)

type Config struct {
	u            *url.URL
	httpClient   HTTPClient
	interceptors []Interceptor
}

func NewClientConfig(serverURL string, opts ...Option) (*Config, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	cfg := &Config{u: u}

	cfg.applyOptions(opts)

	return cfg, nil
}

func (cfg *Config) GetClient() *Client {
	return &Client{
		httpClient:  cfg.httpClient,
		interceptor: cfg.chainInterceptors(),
	}
}

func (cfg *Config) GetURL() *url.URL {
	return cfg.u
}

// chainInterceptors chains all interceptors into one.
func (cfg *Config) chainInterceptors() Interceptor {
	if len(cfg.interceptors) == 0 {
		return nil
	}

	if len(cfg.interceptors) == 1 {
		return cfg.interceptors[0]
	}

	chainedInt := func(ctx context.Context, request any, response APIResp, invoker Invoker) error {
		return cfg.interceptors[0](ctx, request, response, getChainInvoker(cfg.interceptors, 0, invoker))
	}

	return chainedInt
}

func getChainInvoker(interceptors []Interceptor, curr int, finalInvoker Invoker) Invoker {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(ctx context.Context, request any, response APIResp) error {
		return interceptors[curr+1](ctx, request, response, getChainInvoker(interceptors, curr+1, finalInvoker))
	}
}

func (cfg *Config) applyOptions(opts []Option) {
	cfg.httpClient = http.DefaultClient

	for _, opt := range opts {
		opt(cfg)
	}
}
