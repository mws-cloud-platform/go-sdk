package examples_test

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/mws/page"
	secretmanagerclient "go.mws.cloud/go-sdk/service/secretmanager/client"
	secretmanagermodel "go.mws.cloud/go-sdk/service/secretmanager/model"
	secretmanagersdk "go.mws.cloud/go-sdk/service/secretmanager/sdk"
)

// This example demonstrates creating, reading, updating, activating,
// deactivating and deleting a secret versions in secret named
// "example-secret-with-versions" inside the `$MWS_PROJECT` project.
func Example_secret_version() {
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
		log.Panicln("create secret client:", err)
	}

	// Create a new secret version client using the provided SDK.
	secretVersionClient, err := secretmanagersdk.NewSecretVersion(ctx, sdk)
	if err != nil {
		log.Panicln("create secret version client:", err)
	}

	// Use example secret name for demonstration purposes.
	const (
		secretName = "example-secret-with-versions"
	)

	// Create a new secret with version.
	secret, err := secretClient.CreateSecretWithSecretVersion(ctx, secretmanagerclient.CreateSecretWithSecretVersionRequest{
		Name: secretName,
		Body: secretmanagermodel.CreateSecretWithSecretVersionRequest{
			SecretVersionSpecRequest: secretmanagermodel.SecretVersionSpecRequest{
				Data: secretmanagermodel.SecretVersionDataSpec{"keyA": "A1", "keyB": "B1"},
			},
		},
	})
	if err != nil {
		log.Panicln("create secret with version:", err)
	}
	fmt.Printf("secret %s with version %s created\n", secret.GetMetadata().GetId().ResourceName(), secret.Status.CurrentSecretVersion.GetVersion())

	firstVersionRef := secret.Status.CurrentSecretVersion

	// Add a new version to secret.
	version, err := secretVersionClient.AddSecretVersion(ctx, secretmanagerclient.AddSecretVersionRequest{
		Name: secretName,
		Body: secretmanagermodel.AddSecretVersionRequest{
			Spec: &secretmanagermodel.SecretVersionSpecRequest{
				Data: secretmanagermodel.SecretVersionDataSpec{"keyA": "A2", "keyB": "B2"},
			},
		},
	})
	if err != nil {
		log.Panicln("add version to secret:", err)
	}
	fmt.Println("version added:", version.GetMetadata().GetId().ResourceName())

	secondVersionName := string(version.GetMetadata().GetId().ResourceName())

	// Get version info.
	optVersion, err := secretVersionClient.GetSecretVersion(ctx, secretmanagerclient.GetSecretVersionRequest{
		Name:    secretName,
		Version: string(version.GetMetadata().GetId().ResourceName()),
	})
	if err != nil {
		log.Panicln("get secret version:", err)
	}
	fmt.Println("version info received:", optVersion.GetMetadata().GetId().ResourceName())

	// Get version data.
	data, err := secretVersionClient.GetData(ctx, secretmanagerclient.GetDataRequest{
		Name:    secretName,
		Version: string(version.GetMetadata().GetId().ResourceName()),
	})
	if err != nil {
		log.Panicln("get secret version data:", err)
	}
	fmt.Println("version data received:", data)

	// List secret versions.
	resourceNames := make([]string, 0)
	pager := page.NewPager(secretmanagerclient.ListSecretVersionsRequest{
		Name:     secretName,
		PageSize: ptr.Get(10),
	}, secretVersionClient.ListSecretVersions)
	for resource, err := range pager.All(ctx) {
		if err != nil {
			log.Panicln("list secret versions:", err)
		}
		resourceNames = append(resourceNames, string(resource.GetMetadata().GetId().ResourceName()))
	}
	fmt.Println("secret versions listed:", strings.Join(resourceNames, ", "))

	// Deactivate second version.
	_, err = secretVersionClient.UpdateSecretVersion(ctx, secretmanagerclient.UpdateSecretVersionRequest{
		Name:    secretName,
		Version: secondVersionName,
		Body: (&secretmanagermodel.SecretVersionRequest{
			Spec: secretmanagermodel.SecretVersionSpecRequest{
				Active: ptr.Get(false),
			},
		}).AsUpdateModel(),
	})
	if err != nil {
		log.Panicln("deactivate second version:", err)
	}
	fmt.Println("second version deactivated")

	// Activate second version.
	_, err = secretVersionClient.UpdateSecretVersion(ctx, secretmanagerclient.UpdateSecretVersionRequest{
		Name:    secretName,
		Version: secondVersionName,
		Body: (&secretmanagermodel.SecretVersionRequest{
			Spec: secretmanagermodel.SecretVersionSpecRequest{
				Active: ptr.Get(true),
			},
		}).AsUpdateModel(),
	})
	if err != nil {
		log.Panicln("activate second version:", err)
	}
	fmt.Println("second version activated")

	// Set current secret version.
	secret, err = secretClient.UpdateSecret(ctx, secretmanagerclient.UpdateSecretRequest{
		Name: secretName,
		Body: (&secretmanagermodel.SecretRequest{
			Spec: secretmanagermodel.SecretSpecRequest{
				CurrentSecretVersion: firstVersionRef,
			},
		}).AsUpdateModel(),
	})
	if err != nil {
		log.Panicln("set current secret version:", err)
	}
	fmt.Println("set current secret version:", secret.Status.GetCurrentSecretVersion().GetVersion())

	// Delete second secret version.
	err = secretVersionClient.DeleteSecretVersion(ctx, secretmanagerclient.DeleteSecretVersionRequest{
		Name:    secretName,
		Version: secondVersionName,
	})
	if err != nil {
		log.Panicln("delete second secret version:", err)
	}
	fmt.Println("second secret version deleted")

	// And finally, delete the secret to clean up after the example run.
	err = secretClient.DeleteSecret(ctx, secretmanagerclient.DeleteSecretRequest{
		Name: secretName,
	})
	if err != nil {
		log.Panicln("delete secret with versions:", err)
	}
	fmt.Println("secret with versions deleted")
}
