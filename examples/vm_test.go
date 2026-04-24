package examples_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/internal/client"
	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/duration"
	computeclient "go.mws.cloud/go-sdk/service/compute/client"
	computemodel "go.mws.cloud/go-sdk/service/compute/model"
	computesdk "go.mws.cloud/go-sdk/service/compute/sdk"
	computeref "go.mws.cloud/go-sdk/service/resources/references/compute"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"
)

// This example demonstrates creating, reading, updating, and deleting a virtual
// machine named "example-virtual-machine" inside the `$MWS_PROJECT` project.
func Example_setupVirtualMachine() {
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
	virtualMachineClient, err := computesdk.NewVirtualMachine(ctx, sdk)
	if err != nil {
		log.Panicln("create virtual machine client:", err)
	}

	// Use example names for demonstration purposes.
	const (
		virtualMachineName  = "example-virtual-machine"
		networkName         = "example-network"
		subnetName          = "example-subnet"
		externalAddressName = "example-external-address"
	)

	// Create network and subnet required for virtual machine. Clean them up
	// after the example run.
	deleteNetwork := createNetwork(ctx, sdk, networkName)
	defer deleteNetwork()
	subnetRef, deleteSubnet := createSubnet(ctx, sdk, subnetName, networkName)
	defer deleteSubnet()
	externalAddressRef, deleteExternalAddress := createExternalAddress(ctx, sdk, externalAddressName)
	defer deleteExternalAddress()

	// Create a new virtual machine with 2 CPU and 8 GB RAM in the
	// "ru-central1-a" availability zone. Virtual machine would have a boot disk
	// with ubuntu image, address in the "example-subnet" subnet and external address
	// "example-external-address". External address will be selected automatically.
	// Cleanup virtual machine after the example run.
	createVM(ctx, virtualMachineClient, virtualMachineName, subnetRef, externalAddressRef)
	defer deleteVM(ctx, virtualMachineClient, virtualMachineName)

	// Update virtual machine type to 2 CPU and 16 GB RAM.
	updateVM(ctx, virtualMachineClient, virtualMachineName)

	// Get virtual machine by name.
	getVM(ctx, virtualMachineClient, virtualMachineName)
}

func createVM(ctx context.Context, virtualMachineClient *computesdk.VirtualMachine, virtualMachineName string, subnetRef vpcref.SubnetRef, externalAddressRef vpcref.ExternalAddressRef) {
	virtualMachine, err := virtualMachineClient.CreateVirtualMachine(ctx, computeclient.UpsertVirtualMachineRequest{
		VirtualMachine: virtualMachineName,
		Body: computemodel.VirtualMachineRequest{
			Spec: computemodel.VirtualMachineSpecRequest{
				VmType: computeref.NewVmTypeRef("gen-2-8"),
				Zone:   "ru-central1-a",
				Hardware: &computemodel.HardwareSpecRequest{
					Power:                   ptr.Get(computemodel.HardwareSpecPowerRequest_OFF),
					GracefulShutdownTimeout: ptr.Get(duration.NewFromTimeDuration(90 * time.Second)),
				},
				Storage: computemodel.StorageSpecRequest{
					Disks: []computemodel.StorageDiskSpecOrRefWithAttachmentsRequest{
						{
							Name: "boot",
							Boot: ptr.Get(true),
							Disk: computemodel.StorageDiskSpecOrRefRequest{
								Spec: &computemodel.StorageDiskSpecRequest{
									DiskType: ptr.Get(computeref.NewDiskTypeRef("nbs-pl2")),
									Iops:     ptr.Get(computemodel.Iops(1000)),
									Size:     ptr.Get(bytesize.MustParseString("10 GB")),
									Source: &computemodel.StorageDiskSpecSourceRequest{
										Image: ptr.Get(computeref.NewImageRef(
											"mws-ubuntu",
											"mws-ubuntu-2404-lts-v20260324",
										)),
									},
								},
							},
						},
					},
				},
				Network: computemodel.NetworkSpecRequest{
					NetworkInterfaces: []computemodel.NetworkInterfaceSpecRequest{
						{
							Name:    virtualMachineName + "-network-interface-primary",
							Primary: ptr.Get(true),
							Addresses: []computemodel.AddressSpecOrRefWithAttachmentsRequest{
								{
									Address: computemodel.AddressSpecOrRefRequest{
										Spec: &computemodel.AddressSpecRequest{
											Subnet: subnetRef,
										},
									},
									OneToOneNat: &computemodel.ComputeOneToOneNatSpecRequest{
										External: computemodel.ComputeOneToOneNatSpecExternalRequest{
											Address: computemodel.OneToOneNatAddressSpecOrRefRequest{
												Ref: &externalAddressRef,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, computeclient.WithWait())
	if err != nil {
		log.Panicln("create virtual machine:", err)
	}
	fmt.Println("virtual machine created:", ptr.Get(virtualMachine.GetMetadata().GetId()).ResourceName())
}

func getVM(ctx context.Context, virtualMachineClient *computesdk.VirtualMachine, virtualMachineName string) {
	virtualMachine, err := virtualMachineClient.GetVirtualMachine(ctx, computeclient.GetVirtualMachineRequest{
		VirtualMachine: virtualMachineName,
	})
	if err != nil {
		log.Panicln("get virtual machine:", err)
	}
	fmt.Println("virtual machine received:", ptr.Get(virtualMachine.GetMetadata().GetId()).ResourceName())
}

func updateVM(ctx context.Context, virtualMachineClient *computesdk.VirtualMachine, virtualMachineName string) {
	virtualMachine, err := virtualMachineClient.UpdateVirtualMachine(ctx, computeclient.UpdateVirtualMachineRequest{
		VirtualMachine: virtualMachineName,
		Body: computemodel.UpdateVirtualMachineRequest{
			Spec: client.NewOptional(computemodel.UpdateVirtualMachineSpecRequest{
				VmType: client.NewOptional(computeref.NewVmTypeRef("gen-2-16")),
			}),
		},
	}, computeclient.WithWait())
	if err != nil {
		log.Panicln("update virtual machine:", err)
	}
	fmt.Println("virtual machine updated:", ptr.Get(virtualMachine.GetMetadata().GetId()).ResourceName())
}

func deleteVM(ctx context.Context, virtualMachineClient *computesdk.VirtualMachine, virtualMachineName string) {
	err := virtualMachineClient.DeleteVirtualMachine(ctx, computeclient.DeleteVirtualMachineRequest{
		VirtualMachine: virtualMachineName,
	}, computeclient.WithWait())
	if err != nil {
		log.Panicln("delete virtual machine:", err)
	}
	fmt.Println("virtual machine deleted")
}

func createNetwork(ctx context.Context, sdk *mws.SDK, networkName string) func() {
	networkClient, err := vpcsdk.NewNetwork(ctx, sdk)
	if err != nil {
		log.Panicln("create network client:", err)
	}

	network, err := networkClient.CreateNetwork(ctx, vpcclient.UpsertNetworkRequest{
		Network: networkName,
		Body: vpcmodel.NetworkRequest{
			Spec: vpcmodel.VpcNetworkSpecRequest{
				InternetAccess: ptr.Get(true),
				Mtu:            ptr.Get(int32(1500)),
			},
		},
	}, vpcclient.WithWait())
	if err != nil {
		log.Panicln("create network:", err)
	}
	fmt.Println("network created:", network.GetMetadata().GetId().ResourceName())

	return func() {
		err = networkClient.DeleteNetwork(ctx, vpcclient.DeleteNetworkRequest{
			Network: networkName,
		}, vpcclient.WithWait())
		if err != nil {
			log.Panicln("delete network:", err)
		}
		fmt.Println("network deleted:", network.GetMetadata().GetId().ResourceName())
	}
}

func createSubnet(ctx context.Context, sdk *mws.SDK, subnetName, networkName string) (vpcref.SubnetRef, func()) {
	subnetClient, err := vpcsdk.NewSubnet(ctx, sdk)
	if err != nil {
		log.Panicln("create subnet client:", err)
	}

	subnet, err := subnetClient.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: networkName,
		Subnet:  subnetName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("192.168.0.0/16"),
			},
		},
	}, vpcclient.WithWait())
	if err != nil {
		log.Panicln("create subnet:", err)
	}
	fmt.Println("subnet created:", subnet.GetMetadata().GetId().ResourceName())

	deleteSubnet := func() {
		err = subnetClient.DeleteSubnet(ctx, vpcclient.DeleteSubnetRequest{
			Network: networkName,
			Subnet:  subnetName,
		}, vpcclient.WithWait())
		if err != nil {
			log.Panicln("delete subnet:", err)
		}
		fmt.Println("subnet deleted:", subnet.GetMetadata().GetId().ResourceName())
	}

	subnetID, err := vpcref.NewSubnetIDFromAnyID(*subnet.GetMetadata().GetId())
	if err != nil {
		log.Panicln("get subnet id:", err)
	}

	return vpcref.NewSubnetRef(subnetID.GetProject(), subnetID.GetNetwork(), subnetID.GetSubnet()), deleteSubnet
}

func createExternalAddress(ctx context.Context, sdk *mws.SDK, externalAddressName string) (vpcref.ExternalAddressRef, func()) {
	externalAddressClient, err := vpcsdk.NewExternalAddress(ctx, sdk)
	if err != nil {
		log.Panicln("create external address client:", err)
	}

	externalAddress, err := externalAddressClient.CreateExternalAddress(ctx, vpcclient.UpsertExternalAddressRequest{
		ExternalAddress: externalAddressName,
		Body: &vpcmodel.ExternalAddressRequest{
			Spec: vpcmodel.VpcExternalAddressSpecRequest{},
		},
	}, vpcclient.WithWait())
	if err != nil {
		log.Panicln("create external address:", err)
	}
	fmt.Println("external address created:", externalAddress.GetMetadata().GetId().ResourceName())

	deleteExternalAddress := func() {
		err = externalAddressClient.DeleteExternalAddress(ctx, vpcclient.DeleteExternalAddressRequest{
			ExternalAddress: externalAddressName,
		}, vpcclient.WithWait())
		if err != nil {
			log.Panicln("delete external address:", err)
		}
		fmt.Println("external address deleted:", externalAddress.GetMetadata().GetId().ResourceName())
	}

	externalAddressID, err := vpcref.NewExternalAddressIDFromAnyID(*externalAddress.GetMetadata().GetId())
	if err != nil {
		log.Panicln("get external address id:", err)
	}

	return vpcref.NewExternalAddressRef(externalAddressID.GetProject(), externalAddressID.GetExternalAddress()), deleteExternalAddress
}
