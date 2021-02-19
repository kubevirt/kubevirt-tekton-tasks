# Unit Tester Pipeline

Clones kubevirt-tekton-tasks repository and executes unit tests in a VM and then deletes the VM at the end.


## Prerequisites

- KubeVirt `v0.36.0`
- Tekton Pipelines `v0.19.0`

## Pipeline Description

```
  create-vm-from-manifest --- execute-in-vm --- cleanup-vm (finally)
```

1. firstly fedora-unit-tester-private-key and fedora-unit-tester-public-key secrets are created with private and public keys.
2. create-vm-from-manifest task creates a VM called fedora-unit-tester with a public key attached
3. execute-in-vm starts a VM and makes SSH connection to it
    - installs dependencies
    - clones kubevirt-tekton-tasks repository
    - runs unit tests
4. cleanup-vm attempts to connect to the VM
    - prints unit tests results or `failure` if no test results found
    - deletes the VM


## How to run

```bash
kubectl create -f unit-tester-pipeline.yaml
kubectl create -f unit-tester-pipelinerun.yaml
```
