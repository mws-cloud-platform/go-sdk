package examples_test

import (
	"context"
	"fmt"
	"log"

	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/mws"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"
)

// This example demonstrates creating, reading, updating, and deleting a network
// named "example-network" inside the `$MWS_PROJECT` project.
func Example_network() {
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

	// Create a new network client using the provided SDK.
	networkClient, err := vpcsdk.NewNetwork(ctx, sdk)
	if err != nil {
		log.Panicln("create client:", err)
	}

	// Use example network name for demonstration purposes.
	const networkName = "example-network"

	// Create a new network with enabled internet access and MTU set to 1500.
	network, err := networkClient.CreateNetwork(ctx, vpcclient.UpsertNetworkRequest{
		Network: networkName,
		Body: vpcmodel.NetworkRequest{
			Spec: vpcmodel.VpcNetworkSpecRequest{
				InternetAccess: ptr.Get(true),
				Mtu:            ptr.Get(int32(1500)),
			},
		},
	})
	if err != nil {
		log.Panicln("create network:", err)
	}
	fmt.Println("network created:", network.GetMetadata().GetId().ResourceName())

	// Disable internet access for the network.
	network, err = networkClient.UpdateNetwork(ctx, vpcclient.UpdateNetworkRequest{
		Network: networkName,
		Body: (&vpcmodel.NetworkRequest{
			Spec: vpcmodel.VpcNetworkSpecRequest{
				InternetAccess: ptr.Get(false),
			},
		}).AsUpdateModel(),
	})
	if err != nil {
		log.Panicln("update network:", err)
	}
	fmt.Println("network updated:", network.GetMetadata().GetId().ResourceName())

	// Get network by name.
	network, err = networkClient.GetNetwork(ctx, vpcclient.GetNetworkRequest{
		Network: networkName,
	})
	if err != nil {
		log.Panicln("get network:", err)
	}
	fmt.Println("network received:", network.GetMetadata().GetId().ResourceName(), "with MTU:", network.GetStatus().GetMtuOr(0))

	// And finally, delete the network to clean up after the example run.
	err = networkClient.DeleteNetwork(ctx, vpcclient.DeleteNetworkRequest{
		Network: networkName,
	})
	if err != nil {
		log.Panicln("delete network:", err)
	}
	fmt.Println("network deleted")
}
