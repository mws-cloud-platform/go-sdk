package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserAgentInjector(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		userAgent string
	}{
		{
			name:      "filled user agent",
			userAgent: "user-agent",
		},
		{
			name:      "no user agent",
			userAgent: "",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			rt := &mockRoundTripper{fn: func(req *http.Request, _ int) (*http.Response, error) {
				require.Equal(t, testCase.userAgent, req.Header.Get(HeaderUserAgent))
				return nil, nil
			}}
			injector := UserAgentInjector(testCase.userAgent)(rt)

			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

			_, err := injector.RoundTrip(req) //nolint:bodyclose // response has no body
			require.NoError(t, err)
			require.Equal(t, 1, rt.calls)
		})
	}
}
