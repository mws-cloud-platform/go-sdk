package client

type Option func(*Config)

func WithHTTPClient(client HTTPClient) Option {
	return func(o *Config) {
		o.httpClient = client
	}
}

// WithChainInterceptor returns an Option that specifies the chained
// interceptor for unary RPCs. The first interceptor will be the outermost,
// while the last interceptor will be the innermost wrapper around the real call.
func WithChainInterceptor(interceptors ...Interceptor) Option {
	return func(o *Config) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}
