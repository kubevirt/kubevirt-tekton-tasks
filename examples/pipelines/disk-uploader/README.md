# Disk Uploader Pipeline

This pipeline will deploy a new Virutal Machine (VM) - disk uploader - and then the exported Virutal Machine (VM).

## Prerequisites

- KubeVirt `v1.3.0`
- Tekton Pipelines `v0.58.0`

## Pipeline Description

```
  create-example-vm-from-manifest --- disk-uploader --- create-example-vm-exported-from-manifest
```

1. `create-example-vm-from-manifest` task deploys an example Virtual Machine (VM).
2. `disk-uploader` task deploys the tool to extract disk image and upload it to the container registry.
3. `create-example-vm-exported-from-manifest` task deploys the extracted VM from the container regisry.

## How to run

```bash
kubectl apply -f disk-uploader-serviceaccount.yaml
kubectl apply -f disk-uploader-role.yaml
kubectl apply -f disk-uploader-rolebinding.yaml
kubectl apply -f disk-uploader-secret.yaml
kubectl apply -f disk-uploader-pipeline.yaml
kubectl create -f disk-uploader-pipelinerun.yaml
```
