package examples

import (
	"context"
	"fmt"
	"log"

	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/mws/page"
	computeclient "go.mws.cloud/go-sdk/service/compute/client"
	computesdk "go.mws.cloud/go-sdk/service/compute/sdk"
)

// This example demonstrates how to list virtual machines inside the
// `$MWS_PROJECT` project.
func Example_vmList() {
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

	// Create a new virtual machine client using the provided SDK.
	client, err := computesdk.NewVirtualMachine(ctx, sdk)
	if err != nil {
		log.Panicln("create client:", err)
	}

	// List virtual machines with the page size limit.
	virtualMachines, err := client.ListVirtualMachines(ctx, computeclient.ListVirtualMachinesRequest{
		PageSize: ptr.Get(10),
	})
	if err != nil {
		log.Panicln("list virtual machines:", err)
	}

	// Print the virtual machine identifiers.
	fmt.Println("Virtual Machines:")
	for _, vm := range virtualMachines.GetItems() {
		fmt.Println(vm.GetMetadata().GetId())
	}
}

// This example demonstrates how to list virtual machines inside the
// `$MWS_PROJECT` project using iterators.
func Example_vmListIterators() {
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

	// Create a new virtual machine client using the provided SDK.
	client, err := computesdk.NewVirtualMachine(ctx, sdk)
	if err != nil {
		log.Panicln("create client:", err)
	}

	// Create a pager for listing virtual machines with a page size equal to 5.
	pager := page.NewPager(computeclient.ListVirtualMachinesRequest{
		PageSize: ptr.Get(5),
	}, client.ListVirtualMachines)

	// Iterate over all pages with virtual machines and print their identifiers.
	fmt.Println("Pages:")
	num := 1
	for p, err := range pager.Pages(ctx) {
		if err != nil {
			log.Panicln("get page:", err)
		}
		fmt.Printf("Page %d:\n", num)
		for _, vm := range p {
			fmt.Println(vm.GetMetadata().GetId())
		}
	}

	// Iterate over all virtual machines and print their identifiers.
	fmt.Println("Virtual Machines:")
	for vm, err := range pager.All(ctx) {
		if err != nil {
			log.Panicln("get virtual machine:", err)
		}
		fmt.Println(vm.GetMetadata().GetId())
	}
}
