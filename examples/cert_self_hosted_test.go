package examples_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"time"

	"go.mws.cloud/go-sdk/mws"
	certclient "go.mws.cloud/go-sdk/service/certmanager/client"
	certmodel "go.mws.cloud/go-sdk/service/certmanager/model"
	certsdk "go.mws.cloud/go-sdk/service/certmanager/sdk"
)

// This example demonstrates creating, reading, updating, and deleting a
// self-hosted certificate named "example-cert" inside the `$MWS_PROJECT`
// project.
func Example_certSelfHosted() {
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

	// Generate a self-signed certificate that meet requirements:
	// https://mws.ru/docs/cloud-platform/certmanager/general/cert-overview.html#cert-req.
	certPEM, privateKeyPEM := generateSelfSignedCertificate([]string{"example.ru"})

	cert, err := certClient.CreateCertificate(ctx, certclient.UpsertCertificateRequest{
		Name: certName,
		Body: certmodel.CertificateRequest{
			Spec: certmodel.CertificateSpecRequest{
				SelfManaged: &certmodel.SelfManagedSpecRequest{
					Certificate: certPEM,
					PrivateKey:  privateKeyPEM,
				},
			},
		},
	}, certclient.WithWait())
	if err != nil {
		log.Panicln("create certificate:", err)
	}
	fmt.Println("certificate created:", cert.GetMetadata().Id.ResourceName())

	cert, err = certClient.GetCertificate(ctx, certclient.GetCertificateRequest{
		Name: certName,
	})
	if err != nil {
		log.Panicln("get certificate:", err)
	}
	fmt.Println("certificate retrieved:", cert.GetMetadata().Id.ResourceName())

	newCertPEM, newPrivateKeyPEM := generateSelfSignedCertificate([]string{"example.ru"})
	cert, err = certClient.UpdateCertificate(ctx, certclient.UpdateCertificateRequest{
		Name: certName,
		Body: (&certmodel.CertificateRequest{
			Spec: certmodel.CertificateSpecRequest{
				SelfManaged: &certmodel.SelfManagedSpecRequest{
					Certificate: newCertPEM,
					PrivateKey:  newPrivateKeyPEM,
				},
			},
		}).AsUpdateModel(),
	}, certclient.WithWait())
	if err != nil {
		log.Panicln("update certificate:", err)
	}
	fmt.Println("certificate updated:", cert.GetMetadata().Id.ResourceName())

	content, err := certClient.GetCertificateContent(ctx, certclient.GetCertificateContentRequest{
		Name: certName,
	})
	if err != nil {
		log.Panicln("get certificate content:", err)
	}
	fmt.Println("certificate content retrieved:", len(content.GetCertificate()))

	err = certClient.DeleteCertificate(ctx, certclient.DeleteCertificateRequest{
		Name: certName,
	}, certclient.WithWait())
	if err != nil {
		log.Panicln("delete certificate:", err)
	}
	fmt.Println("certificate deleted")
}

func generateSelfSignedCertificate(domains []string) (string, string) {
	generatedPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panicln("generate key:", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Panicln("generate serial number:", err)
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Your Organization"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              domains,
		IsCA:                  false,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &generatedPrivateKey.PublicKey, generatedPrivateKey)
	if err != nil {
		log.Panicln("create certificate:", err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	privateKey, err := x509.MarshalPKCS8PrivateKey(generatedPrivateKey)
	if err != nil {
		log.Panicln("marshal private key:", err)
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKey,
	})

	return string(certPEM), string(privateKeyPEM)
}
