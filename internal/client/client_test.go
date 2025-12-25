package client_test

import (
	"context"

	"go.mws.cloud/go-sdk/internal/client"
)

func noopInvoke(context.Context, any, client.APIResp) error {
	return nil
}
