# Kubevirt Tekton Tasks

VM specific tasks for Tekton Pipelines

## Usage and Deployment

### Create Virtual Machines

- [create-vm-from-manifest](tasks/create-vm-from-manifest)
- [create-vm-from-template](tasks/create-vm-from-template)

### Create DataVolumes

- [create-datavolume-from-manifest](tasks/create-datavolume-from-manifest)

### Execute Commands in Virtual Machines

- [execute-in-vm](tasks/execute-in-vm): execute commands over SSH
- [cleanup-vm](tasks/cleanup-vm): execute commands and/or stop/delete VMs


## Development Guide

See [Getting Started](docs/getting-started.md) for the environment setup and development workflow.
