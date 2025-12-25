package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"

	commonclient "go.mws.cloud/go-sdk/internal/client"
	"go.mws.cloud/go-sdk/internal/client/interceptors/auth"
	"go.mws.cloud/go-sdk/mws/credentials"
)

func TestAuth(t *testing.T) {
	for _, v := range []struct {
		Name          string
		Provider      credentials.Provider
		RequestBefore any
		RequestAfter  any
		Error         error
	}{
		{
			Name:          "authorized request",
			Provider:      credentials.StaticProvider(credentials.Credentials{AccessToken: "other token"}),
			RequestBefore: &requestWithAuthorization{},
			RequestAfter:  &requestWithAuthorization{token: "Bearer other token"},
		},
		{
			Name: "provider error",
			Provider: credentials.ProviderFunc(func(context.Context) (credentials.Credentials, error) {
				return credentials.Credentials{}, errProvider
			}),
			RequestBefore: &requestWithAuthorization{},
			Error:         errProvider,
		},
		{
			Name:          "nil provider",
			RequestBefore: &requestWithAuthorization{},
			RequestAfter:  &requestWithAuthorization{},
		},
		{
			Name:          "anonymous provider",
			Provider:      credentials.AnonymousProvider(),
			RequestBefore: &requestWithAuthorization{},
			RequestAfter:  &requestWithAuthorization{},
		},
		{
			Name:          "unauthorized request",
			RequestBefore: 42,
			RequestAfter:  42,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			var after any
			interceptor := auth.New(v.Provider)
			err := interceptor(t.Context(), v.RequestBefore, nil, func(_ context.Context, request any, _ commonclient.APIResp) error {
				after = request
				return nil
			})
			if v.Error != nil {
				require.ErrorIs(t, err, v.Error)
				return
			}

			require.NoError(t, err)
			require.Equal(t, v.RequestAfter, after)
		})
	}
}

const errProvider = consterr.Error("provider")

type requestWithAuthorization struct {
	token string
}

func (r *requestWithAuthorization) SetAuthorization(token string) {
	r.token = token
}
