package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/internal/client"
)

func TestRequestIDInjector(t *testing.T) {
	requestID := "request-id"
	ctx := client.WithRequestID(t.Context(), requestID)

	rt := &mockRoundTripper{fn: func(req *http.Request, _ int) (*http.Response, error) {
		require.Equal(t, requestID, req.Header.Get(HeaderRequestIDKey))
		return nil, nil
	}}
	injector := RequestIDInjector(rt)

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	req = req.WithContext(ctx)

	_, err := injector.RoundTrip(req) //nolint:bodyclose // response has no body
	require.NoError(t, err)
	require.Equal(t, 1, rt.calls)
}
