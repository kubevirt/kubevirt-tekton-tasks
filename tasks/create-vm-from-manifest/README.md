# Create VirtualMachine from Template Task

This task creates a VirtualMachine from different sources

### Installation

Install the `create-vm-from-manifest` task

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-vm-from-manifest/manifests/create-vm-from-manifest.yaml
```

Install one of the following rbac permissions to the active namespace
- Permissions for using templates/VMs in active namespace
  ```bash
  kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-vm-from-manifest/manifests/create-vm-from-manifest-namespace-rbac.yaml
  ```
- Permissions for using templates/VMs in the cluster
  ```bash
  TARGET_NAMESPACE="$(kubectl config view --minify --output 'jsonpath={..namespace}')"
  wget -qO - https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-vm-from-manifest/manifests/create-vm-from-manifest-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
  ```

### Service Account

This task should be run with `create-vm-from-manifest-task` serviceAccount.

### Parameters

- **manifest**: YAML manifest of a VirtualMachine resource to be created.
- **namespace**: Namespace where to create the VM. (defaults to manifest namespace or active namespace)
- **dataVolumes**: Add DVs to VM Volumes.
- **ownDataVolumes**: Add DVs to VM Volumes and add VM to DV ownerReferences. These DataVolumes will be deleted once the created VM gets deleted.
- **persistentVolumeClaims**: Add PVCs to VM Volumes.
- **ownPersistentVolumeClaims**: Add PVCs to VM Volumes and add VM to PVC ownerReferences. These PVCs will be deleted once the created VM gets deleted.

### Results

- **name**: The name of a VM that was created.
- **namespace**: The namespace of a VM that was created.

### Usage

Please see [examples](examples)
