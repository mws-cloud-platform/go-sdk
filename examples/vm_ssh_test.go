package examples_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

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
	"golang.org/x/crypto/ssh"
)

// This example demonstrates how to generate an SSH key pair, provision a
// virtual machine with the public key embedded into its cloud-config, connect
// to the virtual machine over SSH with the standard Go SSH client, execute a
// simple command, and clean up the created resources.
func Example_sshIntoVirtualMachine() {
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
		virtualMachineName  = "example-virtual-machine-ssh"
		networkName         = "example-network-ssh"
		subnetName          = "example-subnet-ssh"
		externalAddressName = "example-external-address-ssh"
		firewallRuleName    = "example-firewall-rule-allow-ssh"
	)

	// Generate a new SSH key pair. The public key will be placed into the
	// virtual machine cloud-config, and the private key will be used to connect
	// to the virtual machine over SSH.
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Panicln("generate ssh key:", err)
	}
	sshSigner, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		log.Panicln("create ssh signer:", err)
	}
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		log.Panicln("create ssh public key:", err)
	}
	authorizedKey := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(sshPublicKey)))
	fmt.Println("ssh key generated:", authorizedKey)

	// Build a cloud-config with the generated public key. Cloud-init will apply
	// it on the first virtual machine boot and allow SSH access for the
	// "ubuntu" user.
	cloudConfig := fmt.Sprintf(`#cloud-config
users:
  - name: ubuntu
    groups: sudo
    shell: /bin/bash
    sudo: ALL=(ALL) NOPASSWD:ALL
    ssh_authorized_keys:
    - %s
`, authorizedKey)

	// Create network and subnet required for virtual machine. Clean them up
	// after the example run.
	deleteNetwork := createNetwork(ctx, sdk, networkName)
	defer deleteNetwork()
	subnetRef, deleteSubnet := createSubnet(ctx, sdk, subnetName, networkName)
	defer deleteSubnet()
	externalAddressRef, deleteExternalAddress := createExternalAddress(ctx, sdk, externalAddressName)
	defer deleteExternalAddress()

	// Get the IP address reserved for the external address. The virtual
	// machine will be reachable over SSH via this address.
	externalIPAddress := getExternalIPAddress(ctx, sdk, externalAddressName)
	fmt.Println("external address ip address:", externalIPAddress)

	// Create a new virtual machine with 2 CPU and 8 GB RAM in the
	// "ru-central1-a" availability zone. Virtual machine would have a boot disk
	// with ubuntu image, address in the "example-subnet-ssh" subnet and
	// external address "example-external-address-ssh". The cloud-config with
	// the generated public key is passed via the "user-data" OS attribute.
	// Cleanup virtual machine after the example run.
	createVMWithCloudConfig(ctx, virtualMachineClient, virtualMachineName, subnetRef, externalAddressRef, cloudConfig)
	defer deleteVM(ctx, virtualMachineClient, virtualMachineName)

	// Wait until the virtual machine gets its internal address assigned and
	// the external IP address assigned via one-to-one NAT. The external
	// address must match the IP address reserved earlier.
	vmInternalIPAddress, vmExternalIPAddress := waitVMIPAddresses(ctx, virtualMachineClient, virtualMachineName)
	fmt.Println("virtual machine internal ip address:", vmInternalIPAddress)
	fmt.Println("virtual machine external ip address:", vmExternalIPAddress)

	// Create a firewall rule that allows incoming SSH traffic (tcp:22) from
	// any source to the virtual machine. By default, the network firewall
	// blocks all incoming traffic to the virtual machine, so such a rule is
	// required to connect to the virtual machine over SSH. The firewall
	// evaluates packets after the one-to-one NAT translation, so the rule
	// destination is the internal address of the virtual machine, not the
	// reserved external one. Clean the rule up after the example run.
	deleteFirewallRule := createFirewallRuleAllowSSH(ctx, sdk, networkName, firewallRuleName, vmInternalIPAddress)
	defer deleteFirewallRule()

	// Connect to the virtual machine over SSH with the standard Go SSH client
	// using the generated private key. The virtual machine may need some time
	// to boot and apply the cloud-config, so the connection is retried. Close
	// the connection to the virtual machine after the example run.
	sshClient := dialSSH(ctx, externalIPAddress, sshSigner)
	defer func() {
		if closeErr := sshClient.Close(); closeErr != nil {
			log.Panicln("close ssh connection:", closeErr)
		}
		fmt.Println("ssh connection closed")
	}()

	// Open a new SSH session and execute a simple command. Close the session
	// after the example run.
	session, err := sshClient.NewSession()
	if err != nil {
		log.Panicln("create ssh session:", err)
	}
	defer func() {
		if closeErr := session.Close(); closeErr != nil {
			log.Panicln("close ssh session:", closeErr)
		}
		fmt.Println("ssh session closed")
	}()

	output, err := session.Output(`echo "Hello, MWS Cloud Platform!"`)
	if err != nil {
		log.Panicln("execute command:", err)
	}
	fmt.Println("command output:", strings.TrimSpace(string(output)))
}

func createVMWithCloudConfig(ctx context.Context, virtualMachineClient *computesdk.VirtualMachine, virtualMachineName string, subnetRef vpcref.SubnetRef, externalAddressRef vpcref.ExternalAddressRef, cloudConfig string) {
	virtualMachine, err := virtualMachineClient.CreateVirtualMachine(ctx, computeclient.UpsertVirtualMachineRequest{
		VirtualMachine: virtualMachineName,
		Body: computemodel.VirtualMachineRequest{
			Spec: computemodel.VirtualMachineSpecRequest{
				VmType: computeref.NewMustVmTypeRef("gen-2-8"),
				Zone:   "ru-central1-a",
				Hardware: &computemodel.HardwareSpecRequest{
					Power:                   new(computemodel.HardwareSpecPowerRequest_ON),
					GracefulShutdownTimeout: new(duration.NewFromTimeDuration(90 * time.Second)),
				},
				Os: &computemodel.OsSpecRequest{
					Metadata: &computemodel.OsSpecMetadataRequest{
						Attributes: map[string]string{
							"user-data": cloudConfig,
						},
					},
				},
				Storage: computemodel.StorageSpecRequest{
					Disks: []computemodel.StorageDiskSpecOrRefWithAttachmentsRequest{
						{
							Name: "boot",
							Boot: new(true),
							Disk: computemodel.StorageDiskSpecOrRefRequest{
								Spec: &computemodel.StorageDiskSpecRequest{
									DiskType: new(computeref.NewMustDiskTypeRef("nbs-pl2")),
									Iops:     new(computemodel.Iops(1000)),
									Size:     new(bytesize.MustParseString("10 GB")),
									Source: &computemodel.StorageDiskSpecSourceRequest{
										Image: new(computeref.NewMustImageRef(
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
							Primary: new(true),
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
	fmt.Println("virtual machine created:", new(virtualMachine.GetMetadata().GetId()).ResourceName())
}

// waitVMIPAddresses polls the virtual machine until an internal IP address is
// assigned to one of its network interfaces and an external IP address is
// assigned to it via one-to-one NAT. It returns both addresses. The wait is
// bounded by a context deadline, so in-flight requests are cancelled as well
// when the timeout is exceeded. The deadline is scoped to this helper and
// does not affect the caller context.
func waitVMIPAddresses(ctx context.Context, virtualMachineClient *computesdk.VirtualMachine, virtualMachineName string) (string, string) {
	const (
		pollInterval = 10 * time.Second
		pollTimeout  = 5 * time.Minute
	)

	ctx, cancel := context.WithTimeout(ctx, pollTimeout)
	defer cancel()

	for {
		virtualMachine, err := virtualMachineClient.GetVirtualMachine(ctx, computeclient.GetVirtualMachineRequest{
			VirtualMachine: virtualMachineName,
		})
		if err != nil {
			log.Panicln("get virtual machine:", err)
		}

		var internalIPAddress, externalIPAddress string
		for _, networkInterface := range virtualMachine.GetStatus().GetNetwork().NetworkInterfaces {
			for _, address := range networkInterface.GetAddresses() {
				if internalIPAddress == "" {
					if ipAddress := address.GetIpAddress(); ipAddress != nil {
						internalIPAddress = ipAddress.String()
					}
				}
				if externalIPAddress == "" {
					if ipAddress := address.GetOneToOneNat().GetExternal().GetIpAddress(); ipAddress != nil {
						externalIPAddress = ipAddress.String()
					}
				}
			}
		}

		if internalIPAddress != "" && externalIPAddress != "" {
			return internalIPAddress, externalIPAddress
		}

		select {
		case <-ctx.Done():
			log.Panicln("wait virtual machine ip addresses:", ctx.Err())
		case <-time.After(pollInterval):
		}
	}
}

// dialSSH connects to the virtual machine over SSH with the standard Go SSH
// client. The virtual machine may need some time to boot and apply the
// cloud-config, so the connection is retried until the deadline is exceeded.
// Note that ssh.Dial is not context-aware, so the context deadline only
// bounds the retry loop; each dial attempt is bounded by the ssh config
// timeout.
func dialSSH(ctx context.Context, address string, signer ssh.Signer) *ssh.Client {
	const (
		dialTimeout  = 10 * time.Second
		retryTimeout = 5 * time.Minute
	)

	sshConfig := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout: dialTimeout,
		// In this example the host key of the virtual machine is not verified.
		// In production code, verify the host key instead.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	ctx, cancel := context.WithTimeout(ctx, retryTimeout)
	defer cancel()

	for {
		client, err := ssh.Dial("tcp", net.JoinHostPort(address, "22"), sshConfig)
		if err == nil {
			fmt.Println("ssh connection established")
			return client
		}

		select {
		case <-ctx.Done():
			// Report the last dial error along with the deadline: it usually
			// hints at the actual problem.
			log.Panicln("ssh dial:", ctx.Err(), "; last error:", err)
		case <-time.After(dialTimeout):
		}
	}
}

// getExternalIPAddress returns the IP address reserved for the external
// address resource.
func getExternalIPAddress(ctx context.Context, sdk *mws.SDK, externalAddressName string) string {
	externalAddressClient, err := vpcsdk.NewExternalAddress(ctx, sdk)
	if err != nil {
		log.Panicln("create external address client:", err)
	}

	externalAddress, err := externalAddressClient.GetExternalAddress(ctx, vpcclient.GetExternalAddressRequest{
		ExternalAddress: externalAddressName,
	})
	if err != nil {
		log.Panicln("get external address:", err)
	}

	ipAddress := externalAddress.GetStatus().GetIpAddress()
	if ipAddress == nil {
		log.Panicln("external address has no ip address assigned")
	}
	return ipAddress.String()
}

// createFirewallRuleAllowSSH creates a firewall rule that allows incoming SSH
// traffic (tcp:22) from any source to the given internal IP address of the
// virtual machine. By default, the network firewall blocks all incoming
// traffic to virtual machines, so such a rule is required to connect to the
// virtual machine over SSH. The firewall evaluates packets after the
// one-to-one NAT translation, so the destination is the internal address of
// the virtual machine, not its external address.
func createFirewallRuleAllowSSH(ctx context.Context, sdk *mws.SDK, networkName, firewallRuleName, internalIPAddress string) func() {
	firewallRuleClient, err := vpcsdk.NewFirewallRule(ctx, sdk)
	if err != nil {
		log.Panicln("create firewall rule client:", err)
	}

	firewallRule, err := firewallRuleClient.CreateFirewallRule(ctx, vpcclient.UpsertFirewallRuleRequest{
		Network:      networkName,
		FirewallRule: firewallRuleName,
		Body: vpcmodel.FirewallRuleRequest{
			Spec: vpcmodel.FirewallRuleSpecRequest{
				Direction: vpcmodel.FirewallRuleSpecDirectionRequest_INGRESS,
				// Priority of the rule. The lower the number, the higher the
				// priority. Allowed values are 1000-64535.
				Priority: new(int32(1000)),
				Action:   vpcmodel.FirewallRuleSpecActionRequest_ALLOW,
				Source: vpcmodel.FirewallRuleSourceRequest{
					Spec: &vpcmodel.FirewallRuleSourceSpecRequest{
						Cidrs: []cidraddress.CIDR4Address{
							cidraddress.MustParseCIDR4AddressString("0.0.0.0/0"),
						},
					},
				},
				Destination: vpcmodel.FirewallRuleDestinationRequest{
					Spec: &vpcmodel.FirewallRuleDestinationSpecRequest{
						Cidrs: []cidraddress.CIDR4Address{
							cidraddress.MustParseCIDR4AddressString(internalIPAddress + "/32"),
						},
					},
				},
				ProtoPorts: []string{"tcp:22"},
			},
		},
	}, vpcclient.WithWait())
	if err != nil {
		log.Panicln("create firewall rule:", err)
	}
	fmt.Println("firewall rule created:", firewallRule.GetMetadata().GetId().ResourceName())

	return func() {
		err = firewallRuleClient.DeleteFirewallRule(ctx, vpcclient.DeleteFirewallRuleRequest{
			Network:      networkName,
			FirewallRule: firewallRuleName,
		}, vpcclient.WithWait())
		if err != nil {
			log.Panicln("delete firewall rule:", err)
		}
		fmt.Println("firewall rule deleted:", firewallRule.GetMetadata().GetId().ResourceName())
	}
}
