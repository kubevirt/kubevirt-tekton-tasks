# Virt-sysprep Updater Pipeline

Prepares a PVC from Fedora Cloud URL source, updates and seals the operating system with virt-sysprep.
Then it creates a VM from this PVC.

## Prerequisites

- Tekton Pipelines `v0.11.0`

## Pipeline Description

```
  modify-data-object --- disk-virt-sysprep --- create-vm-from-manifest
```

1. `modify-data-object` task imports a PVC from Fedora Cloud URL source. The name of the PVC is generated.
2. `disk-virt-sysprep` task runs yum update and seals the PVC image with virt-sysprep.
4. `create-vm-from-manifest` task creates a VM called `virt-sysprep-updated-vm-*` from the prepared PVC.

## How to run

```bash
kubectl apply -f virt-sysprep-updater-pipeline.yaml
kubectl create -f virt-sysprep-updater-pipelinerun.yaml
```
