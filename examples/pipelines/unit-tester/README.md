# Unit Tester Pipeline

Clones kubevirt-tekton-tasks repository and executes unit tests in a VM and then deletes the VM at the end.


## Prerequisites

- KubeVirt `v0.36.0`
- Tekton Pipelines `v0.19.0`

## Pipeline Description

```
  generate-ssh-keys --- create-vm-from-manifest --- execute-in-vm --- cleanup-vm (finally)
```

1. `generate-ssh-keys` task generates unit-tester-client-private-key and unit-tester-client-public-key secrets with private and public keys.
2. `create-vm-from-manifest` task creates a VM called fedora-unit-tester with a public key attached.
3. `execute-in-vm` task starts a VM and makes SSH connection to it.
    - installs dependencies
    - clones kubevirt-tekton-tasks repository
    - runs unit tests
4. `cleanup-vm` task attempts to connect to the VM.
    - prints unit tests results or `failure` if no test results found
    - deletes the VM


## How to run

```bash
kubectl apply -f unit-tester-pipeline.yaml
kubectl create -f unit-tester-pipelinerun.yaml
```
