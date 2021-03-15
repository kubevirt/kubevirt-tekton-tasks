# KubeVirt Tekton Tasks

[Tekton Pipelines](https://github.com/tektoncd/pipeline) are CI/CD-style pipelines for Kubernetes.
This repository provides KubeVirt-specific Tekton tasks, which focus on:

- Creating and managing resources (VMs, DataVolumes)
- Executing commands in VMs
- Manipulating disk images with libguestfs tools

## Deployment

### On Kubernetes

In order to install the KubeVirt Tekton tasks in the active namespace you need to apply [this manifest](https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/manifests/kubernetes/kubevirt-tekton-tasks.yaml).
You have to repeat this for every namespace in which you'd like to run the tasks.

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/manifests/kubernetes/kubevirt-tekton-tasks.yaml
```

Visit [RBAC permissions for running the tasks](docs/tasks-rbac-permissions.md) if the pipeline needs to create/access resources (VMs, PVCs, etc.) in a different namespace other than the one the pipeline runs in.

### On OpenShift

In order to install the KubeVirt Tekton tasks with additional OpenShift-specific tasks in the active namespace you need to apply [this manifest](https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/manifests/openshift/kubevirt-tekton-tasks.yaml).
You have to repeat this for every namespace in which you'd like to run the tasks.

```bash
oc apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/manifests/openshift/kubevirt-tekton-tasks.yaml
```

Visit [RBAC permissions for running the tasks](docs/tasks-rbac-permissions.md) if the pipeline needs to create/access resources (VMs, PVCs, etc.) in a different namespace other than the one the pipeline runs in.


## Usage

#### Create Virtual Machines

- [create-vm-from-manifest](tasks/create-vm-from-manifest)
- [create-vm-from-template](tasks/create-vm-from-template)

#### Create DataVolumes

- [create-datavolume-from-manifest](tasks/create-datavolume-from-manifest)

#### Generate SSH Keys

- [generate-ssh-keys](tasks/generate-ssh-keys)

#### Execute Commands in Virtual Machines

- [execute-in-vm](tasks/execute-in-vm): execute commands over SSH
- [cleanup-vm](tasks/cleanup-vm): execute commands and/or stop/delete VMs

#### Manipulate PVCs with libguestfs tools

- [disk-virt-customize](tasks/disk-virt-customize): execute virt-customize commands in PVCs

## Examples

#### [Unit Tester Pipeline](examples/pipelines/unit-tester) 

Good unit tests are detached from the operating system and can run everywhere.
However, this is not always the case. Your tests may require access to entire operating system, or run as root,
or need a specific kernel.

This example shows how you can run your tests in your VM of choice.
The pipeline creates a VM, connects to it over SSH and runs tests inside it.
It also showcases the `finally` construct.


#### [Server Deployer Pipeline](examples/pipelines/server-deployer)

For complex application server deployments it might be easier to start the server as is in a VM rather than converting it to cloud-native application.

This example shows how you could initialize/modify a PVC and deploy such an application in a VM.

## Development Guide

See [Getting Started](docs/getting-started.md) for the environment setup and development workflow.
