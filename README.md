[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/kubevirt-tekton-tasks)](https://artifacthub.io/packages/search?repo=kubevirt-tekton-tasks) [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/kubevirt-tekton-pipelines)](https://artifacthub.io/packages/search?repo=kubevirt-tekton-pipelines)
# KubeVirt Tekton Tasks

[Tekton Pipelines](https://github.com/tektoncd/pipeline) are CI/CD-style pipelines for Kubernetes.
This repository provides KubeVirt-specific Tekton tasks, which focus on:

- Creating and managing resources (VMs, DataVolumes, DataSources)
- Executing commands in VMs
- Manipulating disk images with libguestfs tools

## Deployment

In order to install the KubeVirt Tekton tasks in the active namespace you need to apply the following manifest.
You have to repeat this for every namespace in which you'd like to run the tasks.

```bash
VERSION=$(curl -s https://api.github.com/repos/kubevirt/kubevirt-tekton-tasks/releases | \
            jq '.[] | select(.prerelease==false) | .tag_name' | sort -V | tail -n1 | tr -d '"')
kubectl apply -f "https://github.com/kubevirt/kubevirt-tekton-tasks/releases/download/${VERSION}/kubevirt-tekton-tasks.yaml"
```

Visit [RBAC permissions for running the tasks](docs/tasks-rbac-permissions.md) if the pipeline needs to create/access resources (VMs, PVCs, etc.) in a different namespace other than the one the pipeline runs in.

## Usage

#### copy-template, modify-vm-template and create-vm-from-template tasks are deprecated and will be removed in future versions. These tasks based on templates will be replaced with create-vm task with enhancements related to instance types.

#### Create Virtual Machines

- [create-vm-from-manifest](release/tasks/create-vm-from-manifest)
- [create-vm-from-template](release/tasks/create-vm-from-template)

#### Utilize Templates

- [copy-template](release/tasks/copy-template)
- [modify-vm-template](release/tasks/modify-vm-template)

#### Modify data objects

- [modify-data-object](release/tasks/modify-data-object)

#### Generate SSH Keys

- [generate-ssh-keys](release/tasks/generate-ssh-keys)

#### Execute Commands in Virtual Machines

- [execute-in-vm](release/tasks/execute-in-vm): execute commands over SSH
- [cleanup-vm](release/tasks/cleanup-vm): execute commands and/or stop/delete VMs

#### Manipulate PVCs with libguestfs tools

- [disk-virt-customize](release/tasks/disk-virt-customize): execute virt-customize commands in PVCs
- [disk-virt-sysprep](release/tasks/disk-virt-sysprep): execute virt-sysprep commands in PVCs

#### Wait for Virtual Machine Instance Status

- [wait-for-vmi-status](release/tasks/wait-for-vmi-status)

#### Modify Windows iso
- [modify-windows-iso-file](release/tasks/modify-windows-iso-file) - modifies windows iso (replaces prompt bootloader with no-promt 
   bootloader) and replaces original iso in PVC with updated one.

## Examples

#### [Windows EFI Installer Pipeline](release/pipelines/windows-efi-installer)

Downloads a Windows ISO file into a PVC and automatically installs Windows 10, 11 or Server 2k22 with EFI enabled by using a custom Answer file into a new base DataVolume.
Supported Windows versions: Windows 10, 11, Server 2k22

#### [Windows Customize Pipeline](release/pipelines/windows-customize)

Applies customizations to an existing Windows 10, 11, Server 2k22 installation by using a custom Answer file and creates a new base DataVolume.
Supported Windows versions: Windows 10, 11, Server 2k22

#### [Unit Tester Pipeline](examples/pipelines/unit-tester) - Unmaintained

Good unit tests are detached from the operating system and can run everywhere.
However, this is not always the case. Your tests may require access to entire operating system, or run as root,
or need a specific kernel.

This example shows how you can run your tests in your VM of choice.
The pipeline creates a VM, connects to it over SSH and runs tests inside it.
It also showcases the `finally` construct.


#### [Server Deployer Pipeline](examples/pipelines/server-deployer) - Unmaintained

For complex application server deployments it might be easier to start the server as is in a VM rather than converting it to cloud-native application.

This example shows how you can initialize/modify a PVC and deploy such application in a VM.

#### [Virt-sysprep Updater Pipeline](examples/pipelines/virt-sysprep-updater) - Unmaintained

Virt-sysprep can be used for preparing VM images which can be then used as base images for other VMs.

This example shows how you can update an operating system and seal VM's image by using virt-customize.
Then, a VM is created from such image.

## Development Guide

See [Getting Started](docs/getting-started.md) for the environment setup and development workflow.
