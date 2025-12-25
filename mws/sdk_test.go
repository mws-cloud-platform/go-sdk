package mws_test

import (
	"context"
	"fmt"
	"log"

	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/mws/credentials"
)

func ExampleSDK_staticCredentials() {
	ctx := context.Background()

	sdk, err := mws.Load(ctx, mws.WithCredentials(
		credentials.StaticProvider(credentials.Credentials{
			AccessToken: "example-token",
		}),
	))
	if err != nil {
		log.Panic(err)
	}
	defer sdk.Close(ctx)

	creds, err := sdk.CredentialsProvider().Provide(ctx)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(creds.AccessToken)
	// Output: example-token
}
