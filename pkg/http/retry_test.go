package http

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/internal/client/interceptors/retry"
)

func TestRetryAttemptInjector(t *testing.T) {
	attempt := 1
	ctx := retry.WithAttempt(t.Context(), attempt)

	rt := &mockRoundTripper{fn: func(req *http.Request, _ int) (*http.Response, error) {
		actual, err := strconv.Atoi(req.Header.Get(HeaderRetryAttemptKey))
		require.NoError(t, err)
		require.Equal(t, attempt, actual)
		return nil, nil
	}}
	injector := RetryAttemptInjector(rt)

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	req = req.WithContext(ctx)

	_, err := injector.RoundTrip(req) //nolint:bodyclose // response has no body
	require.NoError(t, err)
	require.Equal(t, 1, rt.calls)
}
