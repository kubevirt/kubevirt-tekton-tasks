# Server Deployer Pipeline

Prepares a PVC from Fedora Cloud URL source and installs additional dependencies with virt-customize.
Then it creates a VM from this PVC and deploys a flaskr server application in this VM.

## Prerequisites

- KubeVirt `v0.36.0`
- Tekton Pipelines `v0.11.0`

## Pipeline Description

```
  create-data-object --- disk-virt-customize --- create-vm-from-manifest --- execute-in-vm
                                                           |
                                         generate-ssh-keys--
```

1. `create-data-object` task imports a PVC from Fedora Cloud URL source. The name of the PVC is generated.
2. `disk-virt-customize` task runs virt-customize commands on the PVC that install git, vim, pip and flask python framework.
3. `generate-ssh-keys` task generates two secrets with private and public keys.
   The name of the secrets are generated. The task itself runs in parallel with `1.` and `2.` tasks.
4. `create-vm-from-manifest` task creates a VM  called `flasker-vm-*` from the prepared PVC with a public key attached.
5. `execute-in-vm` task starts a VM and makes SSH connection to it.
   - clones flask repository
   - initializes flaskr application
   - deploys the server on port 5000

## How to run

```bash
kubectl apply -f server-deployer-pipeline.yaml
kubectl create -f server-deployer-pipelinerun.yaml
```

## Interact with the deployed app

To expose and interact with the deployed application run the following snippet
and visit `http://localhost:8001/api/v1/namespaces/${VM_NAMESPACE}/services/flaskr/proxy/`.

VM name and namespace can be found in the PipelineRun's results once it has finished.

```bash
virtctl expose vm ${VM_NAME} --name=flaskr --port 5000
kubectl proxy -p 8001
```
