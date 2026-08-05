package examples_test

import (
	"context"
	"fmt"
	"log"

	"go.mws.cloud/go-sdk/mws"
	resourcemanagerclient "go.mws.cloud/go-sdk/service/rm/client"
	resourcemanagersdk "go.mws.cloud/go-sdk/service/rm/sdk"
)

// This example demonstrates how to enable a "compute" service inside the
// `$MWS_PROJECT` project.
func Example_enableService() {
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

	// Create a new EnabledService client using the provided SDK.
	client, err := resourcemanagersdk.NewEnabledService(ctx, sdk)
	if err != nil {
		log.Panicln("init client:", err)
	}

	// Enable "compute" service.
	resp, err := client.EnableService(ctx, resourcemanagerclient.EnableServiceRequest{
		Service: "compute",
	})
	if err != nil {
		log.Panicln("request failed:", err)
	}

	fmt.Println("service enabled:", resp.GetMetadata().GetId().ResourceName())
}
