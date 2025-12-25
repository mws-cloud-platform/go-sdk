# Examples


All SDK usage examples are stored in a single folder, each with a detailed description.

**Important:** Examples create real resources for which charges may apply.

Available examples:
* [CRUD Network](./network_test.go): Demonstrates how to perform Create, Read, Update, and Delete operations on network resources.
* [CRUD Disk](./disk_test.go): Demonstrates how to perform Create, Read, Update, and Delete operations on disk resources.
* [CRUD Virtual Machine](./vm_test.go): Demonstrates how to perform Create, Read, Update, and Delete operations on virtual machine resources.
* [List Virtual Machines](./vm_list_test.go): Demonstrates how to list virtual machines.
* [Snapshot](./snapshot_test.go): Demonstrates how to create a disk, a snapshot of it, and how to create a copy of the disk from that snapshot.

## Environment Setup

* **CLI Setup**: Step-by-step instructions are available [here](https://mws.ru/docs/mws-cli/general/quickstart-mws-cli.html).
* **Environment Variables:**
    - Save the token obtained from the `mws iam create-token` command in the `MWS_TOKEN` variable:
        ```shell
        export MWS_TOKEN="$(mws iam create-token)"
        ```
    - Set your project name in the `MWS_PROJECT` variable:
        ```shell
        export MWS_PROJECT="your-sandbox-project"
        ```
