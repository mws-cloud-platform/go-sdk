package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/internal/client"
)

func TestIdempotencyTokenInjector(t *testing.T) {
	t.Run("inject from context", func(t *testing.T) {
		idempotencyKey := "idempotency-key"
		ctx := client.WithIdempotencyKey(t.Context(), idempotencyKey)

		rt := &mockRoundTripper{fn: func(req *http.Request, _ int) (*http.Response, error) {
			require.Equal(t, idempotencyKey, req.Header.Get(HeaderIdempotencyKey))
			return nil, nil
		}}
		injector := IdempotencyTokenInjector(rt)

		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		req = req.WithContext(ctx)

		_, err := injector.RoundTrip(req) //nolint:bodyclose // response has no body
		require.NoError(t, err)
		require.Equal(t, 1, rt.calls)
	})

	t.Run("already set", func(t *testing.T) {
		idempotencyKey := "idempotency-key"
		ctx := client.WithIdempotencyKey(t.Context(), "idempotency-key-from-context")

		rt := &mockRoundTripper{fn: func(req *http.Request, _ int) (*http.Response, error) {
			require.Equal(t, idempotencyKey, req.Header.Get(HeaderIdempotencyKey))
			return nil, nil
		}}
		injector := IdempotencyTokenInjector(rt)

		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		req.Header.Set(HeaderIdempotencyKey, idempotencyKey)
		req = req.WithContext(ctx)

		_, err := injector.RoundTrip(req) //nolint:bodyclose // response has no body
		require.NoError(t, err)
		require.Equal(t, 1, rt.calls)
	})
}
