package auth

import (
	"context"
	"fmt"

	"go.mws.cloud/go-sdk/internal/client"
	"go.mws.cloud/go-sdk/mws/credentials"
	"go.mws.cloud/go-sdk/pkg/apimodels/sensitive"
)

func New(provider credentials.Provider) client.Interceptor {
	return func(ctx context.Context, request any, response client.APIResp, invoker client.Invoker) error {
		if req, ok := request.(authorizedRequest); ok && provider != nil {
			creds, err := provider.Provide(ctx)
			if err != nil {
				return fmt.Errorf("provide credentials: %w", err)
			}
			if v := creds.AccessToken.Value(); v != "" {
				req.SetAuthorization(sensitive.New("Bearer " + v))
			}
		}

		return invoker(ctx, request, response)
	}
}

type authorizedRequest interface {
	SetAuthorization(sensitive.Sensitive[string])
}
