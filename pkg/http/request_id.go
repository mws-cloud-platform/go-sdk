package http

import (
	"net/http"

	"go.mws.cloud/go-sdk/internal/client"
)

// HeaderRequestIDKey is the key for the request ID in the request headers.
const HeaderRequestIDKey = "x-request-id"

// RequestIDInjector is an HTTP client middleware that retrieves request ID from
// the request context and injects it into request headers.
func RequestIDInjector(rt http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if requestID := client.RequestIDFromContext(req.Context()); requestID != "" {
			req.Header.Set(HeaderRequestIDKey, requestID)
		}
		return rt.RoundTrip(req)
	})
}
