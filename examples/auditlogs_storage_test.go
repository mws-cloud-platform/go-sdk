package examples_test

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/mws/page"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/duration"
	"go.mws.cloud/go-sdk/pkg/optional"
	auditlogsclient "go.mws.cloud/go-sdk/service/auditlogs/client"
	auditlogsmodel "go.mws.cloud/go-sdk/service/auditlogs/model"
	auditlogssdk "go.mws.cloud/go-sdk/service/auditlogs/sdk"
)

// This example demonstrates creating, reading, updating, and deleting an audit
// logs storage named "example-storage" inside the `$MWS_PROJECT` project.
func Example_auditlogs_storage() {
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

	// Create a new audit logs storage client using the provided SDK.
	storageClient, err := auditlogssdk.NewStorage(ctx, sdk)
	if err != nil {
		log.Panicln("create client:", err)
	}

	// Use example storage name for demonstration purposes.
	const storageName = "example-storage"

	// Create a new storage with 10 GB size limit and 30 days retention
	// period.
	storage, err := storageClient.CreateStorage(ctx, auditlogsclient.UpsertStorageRequest{
		StorageName: storageName,
		Body: auditlogsmodel.StorageRequest{
			Spec: auditlogsmodel.AuditLogsStorageSpecRequest{
				Retention: auditlogsmodel.AuditLogsStorageSpecRetentionRequest{
					Period: duration.MustParseString("30 d"),
					Size:   new(bytesize.MustParseString("10 GB")),
				},
			},
		},
	}, auditlogsclient.WithWait())
	if err != nil {
		log.Panicln("create storage:", err)
	}
	fmt.Println("storage created:", storage.GetMetadata().GetId().ResourceName())

	// Increase the retention period to 60 days. Only fields set in the update
	// request body are changed.
	storage, err = storageClient.UpdateStorage(ctx, auditlogsclient.UpdateStorageRequest{
		StorageName: storageName,
		Body: auditlogsmodel.UpdateStorageRequest{
			Spec: optional.NewOptional(auditlogsmodel.UpdateAuditLogsStorageSpecRequest{
				Retention: optional.NewOptional(auditlogsmodel.UpdateAuditLogsStorageSpecRetentionRequest{
					Period: optional.NewOptional(duration.MustParseString("60 d")),
				}),
			}),
		},
	}, auditlogsclient.WithWait())
	if err != nil {
		log.Panicln("update storage:", err)
	}
	fmt.Println("storage updated:", storage.GetMetadata().GetId().ResourceName())

	// Get storage by name.
	storage, err = storageClient.GetStorage(ctx, auditlogsclient.GetStorageRequest{
		StorageName: storageName,
	})
	if err != nil {
		log.Panicln("get storage:", err)
	}
	spec := storage.GetSpec()
	retention := spec.GetRetention()
	fmt.Println("storage received:", storage.GetMetadata().GetId().ResourceName(),
		"with retention period:", retention.GetPeriod().String(),
		"and size:", retention.GetSizeOr(bytesize.ByteSize{}).String())

	// List storages in the project.
	resourceNames := make([]string, 0)
	pager := page.NewPager(auditlogsclient.ListStoragesRequest{
		PageSize: new(10),
	}, storageClient.ListStorages)
	for item, err := range pager.All(ctx) {
		if err != nil {
			log.Panicln("list storages:", err)
		}
		resourceNames = append(resourceNames, string(item.GetMetadata().GetId().ResourceName()))
	}
	fmt.Println("storages listed:", strings.Join(resourceNames, ", "))

	// And finally, delete the storage to clean up after the example run.
	err = storageClient.DeleteStorage(ctx, auditlogsclient.DeleteStorageRequest{
		StorageName: storageName,
	}, auditlogsclient.WithWait())
	if err != nil {
		log.Panicln("delete storage:", err)
	}
	fmt.Println("storage deleted")
}
