# Kubevirt Tekton Tasks

Kubevirt specific tasks for [Tekton Pipelines](https://github.com/tektoncd/pipeline) (CI/CD-style pipelines for k8s).

Tasks focus on:

- creating and managing resources (VMs, DataVolumes)
- executing commands in VM's
- manipulating disk images with virt-customize [WIP]

## Usage and Deployment

#### Create Virtual Machines

- [create-vm-from-manifest](tasks/create-vm-from-manifest)
- [create-vm-from-template](tasks/create-vm-from-template)

#### Create DataVolumes

- [create-datavolume-from-manifest](tasks/create-datavolume-from-manifest)

#### Execute Commands in Virtual Machines

- [execute-in-vm](tasks/execute-in-vm): execute commands over SSH
- [cleanup-vm](tasks/cleanup-vm): execute commands and/or stop/delete VMs

## Examples

- [Unit Tester](examples/pipelines/unit-tester) pipeline creates a VM and runs unit tests over SSH. Also showcases finally construct.

## Development Guide

See [Getting Started](docs/getting-started.md) for the environment setup and development workflow.
