package examples_test

import (
	"context"
	"fmt"
	"log"

	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	computeclient "go.mws.cloud/go-sdk/service/compute/client"
	computemodel "go.mws.cloud/go-sdk/service/compute/model"
	computesdk "go.mws.cloud/go-sdk/service/compute/sdk"
	computeref "go.mws.cloud/go-sdk/service/resources/references/compute"
)

// This example demonstrates how to create a disk, a disk backup of it, and how
// to create a copy of the disk from that disk backup inside the `$MWS_PROJECT`
// project.
func Example_diskBackup() {
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

	// Use example names for demonstration purposes.
	const (
		diskName       = "example-disk"
		diskBackupName = diskName + "-backup"
		diskCopyName   = diskName + "-copy"
	)

	// Create a new disk with 1 GB size, 4096 B block size, nbs-pl2 type, and
	// iops limit of 1000. Clean it up after the example run.
	diskID, deleteDisk := createDisk(ctx, sdk, diskName)
	defer deleteDisk()

	// Create disk backup from the disk. Clean it up after the example run.
	diskBackupID, deleteDiskBackup := createDiskBackup(ctx, sdk, diskBackupName, diskID)
	defer deleteDiskBackup()

	// Create a disk copy from the disk backup. Clean it up after the example run.
	_, deleteDiskCopy := createDiskFromDiskBackup(ctx, sdk, diskCopyName, &diskBackupID)
	defer deleteDiskCopy()
}

func createDisk(ctx context.Context, sdk *mws.SDK, diskName string) (computeref.DiskID, func()) {
	return createDiskFromDiskBackup(ctx, sdk, diskName, nil)
}

func createDiskFromDiskBackup(
	ctx context.Context,
	sdk *mws.SDK,
	diskName string,
	diskBackupID *computeref.DiskBackupID,
) (computeref.DiskID, func()) {
	// Create disk client.
	diskClient, err := computesdk.NewDisk(ctx, sdk)
	if err != nil {
		log.Panicln("create disk client:", err)
	}

	// If disk backup provided use it as a disk source.
	var source *computemodel.DiskSpecSourceRequest
	if diskBackupID != nil {
		source = &computemodel.DiskSpecSourceRequest{
			DiskBackup: new(computeref.NewMustDiskBackupRef(diskBackupID.GetProject(), diskBackupID.GetDiskBackup())),
		}
	}

	// Create a new disk.
	disk, err := diskClient.CreateDisk(ctx, computeclient.UpsertDiskRequest{
		Disk: diskName,
		Body: computemodel.DiskRequest{
			Spec: computemodel.DiskSpecRequest{
				BlockSize: new(bytesize.MustParseString("4096 B")),
				DiskType:  new(computeref.NewMustDiskTypeRef("nbs-pl2")),
				Iops:      new(computemodel.Iops(1000)),
				Size:      new(bytesize.MustParseString("1 GB")),
				Zone:      "ru-central1-a",
				Source:    source,
			},
		},
	}, computeclient.WithWait())
	if err != nil {
		log.Panicln("create disk:", err)
	}
	fmt.Println("disk created:", disk.GetMetadata().GetId().ResourceName())

	deleteDisk := func() {
		err = diskClient.DeleteDisk(ctx, computeclient.DeleteDiskRequest{
			Disk: diskName,
		}, computeclient.WithWait())
		if err != nil {
			log.Panicln("delete disk:", err)
		}
		fmt.Println("disk deleted:", disk.GetMetadata().GetId().ResourceName())
	}

	diskID, err := computeref.NewDiskIDFromAnyID(*disk.GetMetadata().GetId())
	if err != nil {
		log.Panicln("get disk id:", err)
	}

	return diskID, deleteDisk
}

func createDiskBackup(ctx context.Context, sdk *mws.SDK, diskBackupName string, diskID computeref.DiskID) (computeref.DiskBackupID, func()) {
	// Create disk backup client.
	diskBackupClient, err := computesdk.NewDiskBackup(ctx, sdk)
	if err != nil {
		log.Panicln("create disk backup client:", err)
	}

	// Create disk backup for the specified disk.
	diskBackup, err := diskBackupClient.UpsertDiskBackup(ctx, computeclient.UpsertDiskBackupRequest{
		DiskBackup: diskBackupName,
		Body: computemodel.DiskBackupRequest{
			Spec: computemodel.DiskBackupSpecRequest{
				Source: computemodel.DiskBackupSourceRequest{
					Disk: &computemodel.DiskBackupSourceDiskRequest{
						Id: computeref.NewMustDiskRef(diskID.GetProject(), diskID.GetDisk()),
					},
				},
			},
		},
	}, computeclient.WithWait())
	if err != nil {
		log.Panicln("create disk backup:", err)
	}
	fmt.Println("disk backup created:", diskBackup.GetMetadata().GetId().ResourceName())

	deleteDiskBackup := func() {
		err = diskBackupClient.DeleteDiskBackup(ctx, computeclient.DeleteDiskBackupRequest{
			DiskBackup: diskBackupName,
		}, computeclient.WithWait())
		if err != nil {
			log.Panicln("delete disk backup:", err)
		}
		fmt.Println("disk backup deleted:", diskBackup.GetMetadata().GetId().ResourceName())
	}

	diskBackupID, err := computeref.NewDiskBackupIDFromAnyID(*diskBackup.GetMetadata().GetId())
	if err != nil {
		log.Panicln("get disk backup id:", err)
	}

	return diskBackupID, deleteDiskBackup
}
