package http

import "net/http"

// Middleware represens an HTTP client middlware.
type Middleware func(http.RoundTripper) http.RoundTripper

// Chain groups middlwares together in a chain.
func Chain(middlewares ...Middleware) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		for i := len(middlewares) - 1; i >= 0; i-- {
			rt = middlewares[i](rt)
		}
		return rt
	}
}
