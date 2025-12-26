# MWS Cloud Platform Go SDK

[![PkgGoDev](https://pkg.go.dev/badge/go.mws.cloud/go-sdk)](https://pkg.go.dev/go.mws.cloud/go-sdk)
[![Go Report Card](https://goreportcard.com/badge/go.mws.cloud/go-sdk)](https://goreportcard.com/report/go.mws.cloud/go-sdk)
![Last Commit](https://img.shields.io/github/last-commit/mws-cloud-platform/go-sdk)
![Go Version](https://img.shields.io/badge/Go-1.25.4%2B-blue)

MWS Cloud Platform SDK for Go.

> ⚠️ SDK is under active development and may make breaking changes.

- [MWS Cloud Platform Go SDK](#mws-cloud-platform-go-sdk)
	- [Installation](#installation)
	- [Getting Started](#getting-started)
	- [Configuration](#configuration)
		- [Environment Variables](#environment-variables)
		- [Functional Options](#functional-options)
	- [Examples](#examples)
	- [Documentation](#documentation)
	- [Get Help](#get-help)
	- [Creators](#creators)
	
## Installation

```shell
go get go.mws.cloud/go-sdk/mws
```

## Getting Started

To get started, you need to setup project with Go modules and install the MWS Go
SDK dependency. This example demonstrates how to list virtual machines inside
the project (see [runnable example](./examples/vm_list_test.go)).

**Setup Project**

```shell
mkdir my-mws-project
cd my-mws-project
go mod init my-mws-project
```

**Install SDK**

```shell
go get go.mws.cloud/go-sdk/mws
```

**Write Code**

Write the following code to `main.go` file:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	computeclient "go.mws.cloud/go-sdk/service/compute/client"
	computesdk "go.mws.cloud/go-sdk/service/compute/sdk"
	"go.mws.cloud/go-sdk/go/pkg/sdk/mws"
)

func main() {
	ctx := context.Background()

	// Use the default loader to load configuration. It will load configuration
	// from the environment variables and sensible defaults. You can override
	// logic using [mws.LoadSDKOption] functional options. Check the [mws.Load]
	// and [mws.Config] for more details.
	sdk, err := mws.Load(ctx)
	if err != nil {
		log.Fatalln("loading sdk error:", err)
	}
	defer sdk.Close(ctx)

	// Create a new virtual machine client using the provided SDK.
	client, err := computesdk.NewVirtualMachine(ctx, sdk)
	if err != nil {
		log.Panicln("creating client error:", err)
	}

	// List virtual machines with the page size limit.
	virtualMachines, err := client.ListVirtualMachines(ctx, computeclient.ListVirtualMachinesRequest{
		PageSize: ptr.Get(10),
	})
	if err != nil {
		log.Panicln("listing virtual machines error:", err)
	}

	// Print the virtual machine identifiers.
	fmt.Println("Virtual Machines:")
	for _, vm := range virtualMachines.GetItems() {
		fmt.Println(vm.GetMetadata().GetId())
	}
}
```

**Compile and Run**

Before running the code, make sure you have set project and IAM token in the environment variables:

```
export MWS_PROJECT="your-project"
export MWS_TOKEN="$(mws iam create-token)"
```

Run the code:

```shell
go run main.go
```

Output:

```
Virtual Machines:
compute/projects/your-project/virtualMachines/vm-1
compute/projects/your-project/virtualMachines/vm-2
```

## Configuration

MWS Go SDK requires configuration, like credentials and project identifier. You
can provide this information using environment variables and functional options.

### Environment Variables

- `MWS_BASE_ENDPOINT` - MWS Cloud Platform API base endpoint (default: `https://api.mwsapis.ru`).
- `MWS_PROJECT` - Default project identifier.
- `MWS_ZONE` - Default zone identifier (default: `ru-central1-a`).
- `MWS_TOKEN` - IAM token for authentication. If not empty, it will be used in all client requests that require authentication.
- `MWS_SERVICE_ACCOUNT_AUTHORIZED_KEY_PATH` - Path to the service account authorized key file used for authentication. Has no effect if `MWS_TOKEN` is not empty.
- `MWS_TIMEOUT` - Timeout for all client requests (default: `5s`).

### Functional Options

You can also configure SDK using functional options `mws.LoadSDKOption`, for example:

```go
sdk, err := mws.LoadSDK(
	mws.WithDefaultProject("my-project"),
	mws.WithDefaultZone("ru-central1-a"),
	mws.WithTimeout(5 * time.Second),
)
```

Note that functional options have highest priority and overrides behavior based
on the environment variables and configuration defaults.

## Examples

Check more examples in the [examples](./examples) directory.

## Documentation

* [MWS Cloud Platform Documentation](https://mws.ru/docs)

## Get Help

Ask for help using the [MWS Cloud Platform Support Center](https://mws.ru/docs/support/about.html).

## Creators

Created and maintained by [MWS Cloud Platform](https://mws.ru/cloud-platform).
