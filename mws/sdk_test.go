package mws_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"log"

	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/mws/credentials"
	"go.mws.cloud/go-sdk/mws/iam"
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

func ExampleSDK_withServiceAccountAuthorizedKey() {
	ctx := context.Background()

	parsed, err := x509.ParsePKCS8PrivateKey([]byte("<private-key>"))
	if err != nil {
		log.Panic(err)
	}

	privateKey, ok := parsed.(*ecdsa.PrivateKey)
	if !ok {
		log.Panic("parsed key is not an ECDSA private key")
	}

	sdk, err := mws.Load(ctx, mws.WithServiceAccountAuthorizedKey(
		iam.ServiceAccountAuthorizedKey{
			ServiceAccount: iam.ServiceAccount{
				Project: "<project>",
				Name:    "<name>",
			},
			AuthorizedKey: iam.AuthorizedKey{
				Name:       "<name>",
				PrivateKey: privateKey,
				Algorithm:  "ES256",
			},
		},
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
}

func ExampleSDK_serviceAccountAuthorizedKeyFromFile() {
	ctx := context.Background()

	key, err := iam.ServiceAccountAuthorizedKeyFromFile("/path/to/key.json")
	if err != nil {
		log.Panic(err)
	}

	sdk, err := mws.Load(ctx, mws.WithServiceAccountAuthorizedKey(key))
	if err != nil {
		log.Panic(err)
	}
	defer sdk.Close(ctx)

	creds, err := sdk.CredentialsProvider().Provide(ctx)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(creds.AccessToken)
}
