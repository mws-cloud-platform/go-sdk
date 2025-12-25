package http

import (
	"net/http"
	"strconv"

	"go.mws.cloud/go-sdk/internal/client/interceptors/retry"
)

// HeaderRetryAttemptKey is the key for the attempt number in the request headers.
const HeaderRetryAttemptKey = "x-retry-attempt"

// RetryAttemptInjector is an HTTP client middleware that retrieves attempt number from
// the request context and injects it into request headers.
func RetryAttemptInjector(rt http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if attempt := retry.AttemptFromContext(req.Context()); attempt != 0 {
			req.Header.Set(HeaderRetryAttemptKey, strconv.Itoa(attempt))
		}
		return rt.RoundTrip(req)
	})
}
