package http

import (
	"net/http"
)

// HeaderUserAgent is the key for the user-agent in the request headers.
const HeaderUserAgent = "User-Agent"

// UserAgentInjector is an HTTP client middleware that injects user-agent
// into request headers.
func UserAgentInjector(userAgent string) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if userAgent != "" {
				req.Header.Set(HeaderUserAgent, userAgent)
			}
			return rt.RoundTrip(req)
		})
	}
}
