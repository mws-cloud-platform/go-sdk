package http

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Otel is an HTTP client middleware that adds OpenTelemetry instrumentation.
func Otel(opts ...otelhttp.Option) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return otelhttp.NewTransport(rt, opts...)
	}
}
