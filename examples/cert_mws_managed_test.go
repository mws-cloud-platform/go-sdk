package examples_test

import (
	"context"
	"fmt"
	"log"

	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/mws"
	certclient "go.mws.cloud/go-sdk/service/certmanager/client"
	certmodel "go.mws.cloud/go-sdk/service/certmanager/model"
	certsdk "go.mws.cloud/go-sdk/service/certmanager/sdk"
	common "go.mws.cloud/go-sdk/service/common/model"
)

// This example demonstrates creating, reading, updating, and deleting a
// mws-managed certificate named "example-cert" inside the `$MWS_PROJECT`
// project.
func Example_certMWSManaged() {
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

	certClient, err := certsdk.NewCertificate(ctx, sdk)
	if err != nil {
		log.Panicln("create client:", err)
	}

	// Use example certificate name for demonstration purposes.
	const certName = "example-cert"

	// Immediately after creation, the certificate will be assigned the
	// "Verifying" status. After this, you have to delegate domain authority
	// verification.
	//
	// See https://mws.ru/docs/cloud-platform/certmanager/general/lets-encrypt-operations.html#domain-access.
	cert, err := certClient.CreateCertificate(ctx, certclient.UpsertCertificateRequest{
		Name: certName,
		Body: certmodel.CertificateRequest{
			Spec: certmodel.CertificateSpecRequest{
				Managed: &certmodel.CertificateManagedSpecRequest{
					PreferredChallengeType: certmodel.CertificateChallengeType_DNS01,
					Provider:               certmodel.CertificateProvider_LETS_ENCRYPT,
					Domains:                []string{"example.ru"}, // your domains
				},
			},
		},
	})
	if err != nil {
		log.Panicln("create certificate:", err)
	}
	fmt.Println("certificate created:", cert.GetMetadata().Id)

	cert, err = certClient.UpdateCertificate(ctx, certclient.UpdateCertificateRequest{
		Name: certName,
		Body: (&certmodel.CertificateRequest{
			Metadata: &common.CommonTypedResourceMetadataRequest{
				Description: ptr.Get("managed certificate"),
			},
		}).AsUpdateModel(),
	})
	if err != nil {
		log.Panicln("update certificate:", err)
	}
	fmt.Println("certificate updated:", cert.GetMetadata().Id)

	content, err := certClient.GetCertificateContent(ctx, certclient.GetCertificateContentRequest{
		Name: certName,
	})
	if err != nil {
		log.Panicln("get certificate content:", err)
	}
	fmt.Println("certificate content retrieved:", content)

	err = certClient.DeleteCertificate(ctx, certclient.DeleteCertificateRequest{
		Name: certName,
	})
	if err != nil {
		log.Panicln("delete certificate:", err)
	}
	fmt.Println("certificate deleted")
}
