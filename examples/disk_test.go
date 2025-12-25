package examples_test

import (
	"context"
	"fmt"
	"log"

	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	commonclient "go.mws.cloud/go-sdk/internal/client"
	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	computeclient "go.mws.cloud/go-sdk/service/compute/client"
	computemodel "go.mws.cloud/go-sdk/service/compute/model"
	computesdk "go.mws.cloud/go-sdk/service/compute/sdk"
	computeref "go.mws.cloud/go-sdk/service/resources/references/compute"
)

// This example demonstrates creating, reading, updating, and deleting a disk
// named "example-disk" inside the `$MWS_PROJECT` project.
func Example_disk() {
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

	// Create a new disk client using the provided SDK.
	diskClient, err := computesdk.NewDisk(ctx, sdk)
	if err != nil {
		log.Panicln("create client:", err)
	}

	// Use example disk name for demonstration purposes.
	const diskName = "example-disk"

	// Create a new disk with 1 GB size, 4096 B block size, nbs-pl2 type, and
	// iops limit of 1000.
	disk, err := diskClient.CreateDisk(ctx, computeclient.UpsertDiskRequest{
		Disk: diskName,
		Body: computemodel.DiskRequest{
			Spec: computemodel.DiskSpecRequest{
				BlockSize: ptr.Get(bytesize.MustParseString("4096 B")),
				DiskType:  ptr.Get(computeref.NewDiskTypeRef("nbs-pl2")),
				Iops:      ptr.Get(computemodel.Iops(1000)),
				Size:      ptr.Get(bytesize.MustParseString("1 GB")),
				Zone:      "ru-central1-a",
			},
		},
	})
	if err != nil {
		log.Panicln("create disk:", err)
	}
	fmt.Println("disk created:", disk.GetMetadata().GetId().ResourceName())

	// Increase size of the disk to 2 GB. Nota that disk size can not be
	// decreased.
	disk, err = diskClient.UpdateDisk(ctx, computeclient.UpdateDiskRequest{
		Disk: diskName,
		Body: computemodel.UpdateDiskRequest{
			Spec: commonclient.NewOptional(computemodel.UpdateDiskSpecRequest{
				Size: commonclient.NewOptional(bytesize.MustParseString("2 GB")),
			}),
		},
	})
	if err != nil {
		log.Panicln("update disk:", err)
	}
	fmt.Println("disk updated:", disk.GetMetadata().GetId().ResourceName())

	// Get disk by name.
	disk, err = diskClient.GetDisk(ctx, computeclient.GetDiskRequest{
		Disk: diskName,
	})
	if err != nil {
		log.Panicln("get disk:", err)
	}
	fmt.Println("disk received:", disk.GetMetadata().GetId().ResourceName(), "with size:", disk.GetStatus().GetSize().String())

	// And finally, delete the disk to clean up after the example run.
	err = diskClient.DeleteDisk(ctx, computeclient.DeleteDiskRequest{
		Disk: diskName,
	})
	if err != nil {
		log.Panicln("delete disk:", err)
	}
	fmt.Println("disk deleted")
}
