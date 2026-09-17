package examples_test

import (
	"context"
	"fmt"
	"log"

	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/duration"
	"go.mws.cloud/go-sdk/pkg/optional"
	auditlogsclient "go.mws.cloud/go-sdk/service/auditlogs/client"
	auditlogsmodel "go.mws.cloud/go-sdk/service/auditlogs/model"
	auditlogssdk "go.mws.cloud/go-sdk/service/auditlogs/sdk"
	auditlogsref "go.mws.cloud/go-sdk/service/resources/references/auditlogs"
)

// This example demonstrates creating, reading, updating, and deleting an audit
// logs collector named "example-collector" inside the `$MWS_PROJECT` project.
func Example_auditlogs_collector() {
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

	// Create a new audit logs collector client using the provided SDK.
	collectorClient, err := auditlogssdk.NewCollector(ctx, sdk)
	if err != nil {
		log.Panicln("create client:", err)
	}

	// Use example names for demonstration purposes.
	const (
		collectorName = "example-collector"
		storageName   = "example-storage"
	)

	// Create audit logs storage required for the collector destination. Clean
	// it up after the example run.
	storageID, deleteStorage := createStorage(ctx, sdk, storageName)
	defer deleteStorage()

	// Create a new enabled collector that delivers events of the "compute"
	// service to the audit logs storage. When no events sources are specified,
	// events are collected from the project containing the collector.
	collector, err := collectorClient.CreateCollector(ctx, auditlogsclient.UpsertCollectorRequest{
		CollectorName: collectorName,
		Body: auditlogsmodel.CollectorRequest{
			Spec: auditlogsmodel.CollectorSpecRequest{
				Destination: auditlogsmodel.CollectorSpecDestinationRequest{
					Storage: storageID,
				},
				Enabled: true,
				Services: []auditlogsmodel.CollectorServiceRequest{
					{Service: "compute"},
				},
			},
		},
	}, auditlogsclient.WithWait())
	if err != nil {
		log.Panicln("create collector:", err)
	}
	fmt.Println("collector created:", collector.GetMetadata().GetId().ResourceName())

	// Disable the collector. Only fields set in the update request body are
	// changed.
	collector, err = collectorClient.UpdateCollector(ctx, auditlogsclient.UpdateCollectorRequest{
		CollectorName: collectorName,
		Body: auditlogsmodel.UpdateCollectorRequest{
			Spec: optional.NewOptional(auditlogsmodel.UpdateCollectorSpecRequest{
				Enabled: optional.NewOptional(false),
			}),
		},
	}, auditlogsclient.WithWait())
	if err != nil {
		log.Panicln("update collector:", err)
	}
	fmt.Println("collector updated:", collector.GetMetadata().GetId().ResourceName())

	// Get collector by name.
	collector, err = collectorClient.GetCollector(ctx, auditlogsclient.GetCollectorRequest{
		CollectorName: collectorName,
	})
	if err != nil {
		log.Panicln("get collector:", err)
	}
	spec := collector.GetSpec()
	fmt.Println("collector received:", collector.GetMetadata().GetId().ResourceName(), "with enabled:", spec.GetEnabled())

	// List collectors in the project.
	collectors, err := collectorClient.ListCollectors(ctx, auditlogsclient.ListCollectorsRequest{
		PageSize: new(10),
	})
	if err != nil {
		log.Panicln("list collectors:", err)
	}
	for _, item := range collectors.GetItems() {
		fmt.Println("collector listed:", item.GetMetadata().GetId().ResourceName())
	}

	// And finally, delete the collector to clean up after the example run.
	err = collectorClient.DeleteCollector(ctx, auditlogsclient.DeleteCollectorRequest{
		CollectorName: collectorName,
	}, auditlogsclient.WithWait())
	if err != nil {
		log.Panicln("delete collector:", err)
	}
	fmt.Println("collector deleted")
}

func createStorage(ctx context.Context, sdk *mws.SDK, storageName string) (*auditlogsref.StorageID, func()) {
	storageClient, err := auditlogssdk.NewStorage(ctx, sdk)
	if err != nil {
		log.Panicln("create storage client:", err)
	}

	storage, err := storageClient.CreateStorage(ctx, auditlogsclient.UpsertStorageRequest{
		StorageName: storageName,
		Body: auditlogsmodel.StorageRequest{
			Spec: auditlogsmodel.AuditLogsStorageSpecRequest{
				Retention: auditlogsmodel.AuditLogsStorageSpecRetentionRequest{
					Period: duration.MustParseString("30 d"),
				},
			},
		},
	}, auditlogsclient.WithWait())
	if err != nil {
		log.Panicln("create storage:", err)
	}
	fmt.Println("storage created:", storage.GetMetadata().GetId().ResourceName())

	deleteStorage := func() {
		err = storageClient.DeleteStorage(ctx, auditlogsclient.DeleteStorageRequest{
			StorageName: storageName,
		}, auditlogsclient.WithWait())
		if err != nil {
			log.Panicln("delete storage:", err)
		}
		fmt.Println("storage deleted:", storage.GetMetadata().GetId().ResourceName())
	}

	return storage.GetMetadata().GetId(), deleteStorage
}
