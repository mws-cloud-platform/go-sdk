package examples_test

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/mws/page"
	"go.mws.cloud/go-sdk/service/common/model"
	secretmanagerclient "go.mws.cloud/go-sdk/service/secretmanager/client"
	secretmanagermodel "go.mws.cloud/go-sdk/service/secretmanager/model"
	secretmanagersdk "go.mws.cloud/go-sdk/service/secretmanager/sdk"
)

// This example demonstrates creating, reading, updating, activating, deactivating
// and deleting a secret named "example-secret" inside the `$MWS_PROJECT` project.
func Example_secret() {
	ctx := context.Background()

	// Use the default SDK loader. It will load configuration from the
	// environment variables and sensible defaults. You can override logic using
	// [mws.LoadSDKOption] options. Check the [mws.Load] and [mws.Config] for
	// more details.
	sdk, err := mws.Load(ctx)
	if err != nil {
		log.Panicln("load sdk:", err)
	}
	defer sdk.Close(ctx)

	// Create a new secret client using the provided SDK.
	secretClient, err := secretmanagersdk.NewSecret(ctx, sdk)
	if err != nil {
		log.Panicln("create client:", err)
	}

	// Use example secret name for demonstration purposes.
	const (
		secretName = "example-secret"
	)

	// Create a new secret without version.
	secret, err := secretClient.CreateSecret(ctx, secretmanagerclient.UpsertSecretRequest{
		Name: secretName,
		Body: secretmanagermodel.SecretRequest{
			Metadata: &model.CommonTypedResourceMetadataRequest{
				DisplayName: ptr.Get("Display name"),
				Description: ptr.Get("Description of secret"),
			},
			Spec: secretmanagermodel.SecretSpecRequest{},
		},
	}, secretmanagerclient.WithWait())
	if err != nil {
		log.Panicln("create secret:", err)
	}
	fmt.Println("secret created:", secret.GetMetadata().GetId().ResourceName())

	// Update display name and description of secret.
	secret, err = secretClient.UpdateSecret(ctx, secretmanagerclient.UpdateSecretRequest{
		Name: secretName,
		Body: (&secretmanagermodel.SecretRequest{
			Metadata: &model.CommonTypedResourceMetadataRequest{
				DisplayName: ptr.Get("New display name"),
				Description: ptr.Get("New description of secret"),
			},
		}).AsUpdateModel(),
	}, secretmanagerclient.WithWait())
	if err != nil {
		log.Panicln("update secret:", err)
	}
	fmt.Println("secret updated:", secret.GetMetadata().GetId().ResourceName())

	// Activate secret.
	secret, err = secretClient.UpdateSecret(ctx, secretmanagerclient.UpdateSecretRequest{
		Name: secretName,
		Body: (&secretmanagermodel.SecretRequest{
			Spec: secretmanagermodel.SecretSpecRequest{
				Active: ptr.Get(true),
			},
		}).AsUpdateModel(),
	}, secretmanagerclient.WithWait())
	if err != nil {
		log.Panicln("activate secret:", err)
	}
	fmt.Println("secret activated:", secret.GetMetadata().GetId().ResourceName())

	// Deactivate secret.
	secret, err = secretClient.UpdateSecret(ctx, secretmanagerclient.UpdateSecretRequest{
		Name: secretName,
		Body: (&secretmanagermodel.SecretRequest{
			Spec: secretmanagermodel.SecretSpecRequest{
				Active: ptr.Get(false),
			},
		}).AsUpdateModel(),
	}, secretmanagerclient.WithWait())
	if err != nil {
		log.Panicln("deactivate secret:", err)
	}
	fmt.Println("secret deactivated:", secret.GetMetadata().GetId().ResourceName())

	// List secrets.
	resourceNames := make([]string, 0)
	pager := page.NewPager(secretmanagerclient.ListSecretsRequest{
		PageSize: ptr.Get(10),
	}, secretClient.ListSecrets)
	for resource, err := range pager.All(ctx) {
		if err != nil {
			log.Panicln("list secrets:", err)
		}
		resourceNames = append(resourceNames, string(resource.GetMetadata().GetId().ResourceName()))
	}
	fmt.Println("secrets listed:", strings.Join(resourceNames, ", "))

	// Get secret by name.
	secret, err = secretClient.GetSecret(ctx, secretmanagerclient.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		log.Panicln("get secret:", err)
	}
	fmt.Println("secret received:", secret.GetMetadata().GetId().ResourceName())

	// And finally, delete secret to clean up after the example run.
	err = secretClient.DeleteSecret(ctx, secretmanagerclient.DeleteSecretRequest{
		Name: secretName,
	}, secretmanagerclient.WithWait())
	if err != nil {
		log.Panicln("delete secret:", err)
	}
	fmt.Println("secret deleted")
}
