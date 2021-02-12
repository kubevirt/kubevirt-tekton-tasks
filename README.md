# KubeVirt Tekton Tasks

[Tekton Pipelines](https://github.com/tektoncd/pipeline) are CI/CD-style pipelines for Kubernetes.
This repository provides KubeVirt-specific Tekton tasks, which focus on:

- Creating and managing resources (VMs, DataVolumes)
- Executing commands in VMs
- Manipulating disk images with virt-customize [WIP]

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

- Good unit tests are detached from the operating system and can run everywhere.
  However, this is not always the case. Your tests may require access to entire operating system, or run as root,
  or need a specific kernel.
  The [Unit Tester](examples/pipelines/unit-tester) example shows how you can run your tests in your VM of choice.
  The pipeline creates a VM, connects to it over SSH and runs tests inside it.
  It also showcases the `finally` construct.

## Development Guide

See [Getting Started](docs/getting-started.md) for the environment setup and development workflow.
