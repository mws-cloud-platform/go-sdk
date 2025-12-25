package http

import (
	"net/http"

	"go.mws.cloud/go-sdk/internal/client"
)

// HeaderIdempotencyKey is the request header key for the idempotency key.
const HeaderIdempotencyKey = "Idempotency-Key"

// IdempotencyTokenInjector is an HTTP client middleware that retrieves
// idempotency key from the request context and injects it into request headers
// if it's not already set.
func IdempotencyTokenInjector(rt http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		key := client.IdempotencyKeyFromContext(req.Context())
		if key != "" && req.Header.Get(HeaderIdempotencyKey) == "" {
			req.Header.Set(HeaderIdempotencyKey, key)
		}
		return rt.RoundTrip(req)
	})
}
