package examples_test

import (
	"context"
	"fmt"
	"log"

	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	computeclient "go.mws.cloud/go-sdk/service/compute/client"
	computemodel "go.mws.cloud/go-sdk/service/compute/model"
	computesdk "go.mws.cloud/go-sdk/service/compute/sdk"
	computeref "go.mws.cloud/go-sdk/service/resources/references/compute"
)

// This example demonstrates how to create a disk, a snapshot of it, and how to
// create a copy of the disk from that snapshot inside the `$MWS_PROJECT`
// project.
func Example_snapshot() {
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
		diskName     = "example-disk"
		snapshotName = diskName + "-snapshot"
		diskCopyName = diskName + "-copy"
	)

	// Create a new disk with 1 GB size, 4096 B block size, nbs-pl2 type, and
	// iops limit of 1000. Clean it up after the example run.
	diskID, deleteDisk := createDisk(ctx, sdk, diskName)
	defer deleteDisk()

	// Create snapshot from the disk. Clean it up after the example run.
	snapshotID, deleteSnapshot := createSnapshot(ctx, sdk, snapshotName, diskID)
	defer deleteSnapshot()

	// Create a disk copy from the snapshot. Clean it up after the example run.
	_, deleteDiskCopy := createDiskFromSnapshot(ctx, sdk, diskCopyName, &snapshotID)
	defer deleteDiskCopy()
}

func createDisk(ctx context.Context, sdk *mws.SDK, diskName string) (computeref.DiskID, func()) {
	return createDiskFromSnapshot(ctx, sdk, diskName, nil)
}

func createDiskFromSnapshot(ctx context.Context, sdk *mws.SDK, diskName string, snapshotID *computeref.SnapshotID) (computeref.DiskID, func()) {
	// Create disk client.
	diskClient, err := computesdk.NewDisk(ctx, sdk)
	if err != nil {
		log.Panicln("create disk client:", err)
	}

	// If snapshot provided use it as a disk source.
	var source *computemodel.DiskSpecSourceRequest
	if snapshotID != nil {
		source = &computemodel.DiskSpecSourceRequest{
			Snapshot: ptr.Get(computeref.NewSnapshotRef(snapshotID.GetProject(), snapshotID.GetSnapshot())),
		}
	}

	// Create a new disk.
	disk, err := diskClient.CreateDisk(ctx, computeclient.UpsertDiskRequest{
		Disk: diskName,
		Body: computemodel.DiskRequest{
			Spec: computemodel.DiskSpecRequest{
				BlockSize: ptr.Get(bytesize.MustParseString("4096 B")),
				DiskType:  ptr.Get(computeref.NewDiskTypeRef("nbs-pl2")),
				Iops:      ptr.Get(computemodel.Iops(1000)),
				Size:      ptr.Get(bytesize.MustParseString("1 GB")),
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

func createSnapshot(ctx context.Context, sdk *mws.SDK, snapshotName string, diskID computeref.DiskID) (computeref.SnapshotID, func()) {
	// Create snapshot client.
	snapshotClient, err := computesdk.NewSnapshot(ctx, sdk)
	if err != nil {
		log.Panicln("create snapshot client:", err)
	}

	// Create snapshot for the specified disk.
	snapshot, err := snapshotClient.CreateSnapshot(ctx, computeclient.UpsertSnapshotRequest{
		Snapshot: snapshotName,
		Body: computemodel.SnapshotRequest{
			Spec: computemodel.SnapshotSpecRequest{
				Source: computemodel.SnapshotSourceRequest{
					Disk: &computemodel.SnapshotSourceDiskRequest{
						Id: computeref.NewDiskRef(diskID.GetProject(), diskID.GetDisk()),
					},
				},
			},
		},
	}, computeclient.WithWait())
	if err != nil {
		log.Panicln("create snapshot:", err)
	}
	fmt.Println("snapshot created:", snapshot.GetMetadata().GetId().ResourceName())

	deleteSnapshot := func() {
		err = snapshotClient.DeleteSnapshot(ctx, computeclient.DeleteSnapshotRequest{
			Snapshot: snapshotName,
		}, computeclient.WithWait())
		if err != nil {
			log.Panicln("delete snapshot:", err)
		}
		fmt.Println("snapshot deleted:", snapshot.GetMetadata().GetId().ResourceName())
	}

	snapshotID, err := computeref.NewSnapshotIDFromAnyID(*snapshot.GetMetadata().GetId())
	if err != nil {
		log.Panicln("get snapshot id:", err)
	}

	return snapshotID, deleteSnapshot
}
