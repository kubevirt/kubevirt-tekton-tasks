# Create VirtualMachine Task

Flavors of this task create a VirtualMachine from different sources

## `create-vm-from-template`

### Installation

Install the Task

```bash
kubectl apply -f https://raw.githubusercontent.com/suomiy/kubevirt-tekton-tasks/master/tasks/create-vm/manifests/create-vm-from-template.yaml
```

Install one of the following rbac permissions to the active namespace
  - Permissions for templates from active namespace
    ```bash
    kubectl apply -f https://raw.githubusercontent.com/suomiy/kubevirt-tekton-tasks/master/tasks/create-vm/manifests/create-vm-namespace-rbac.yaml
    ```
  - Permissions for all templates from the cluster
    ```bash
    TARGET_NAMESPACE="$(kubectl config current-context | cut -d/ -f1)"
    wget -qO - https://raw.githubusercontent.com/suomiy/kubevirt-tekton-tasks/master/tasks/create-vm/manifests/create-vm-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
    ```

### Parameters

- **templateName**: Name of a template to create VM from
- **templateNamespace**: Namespace of a template to create VM from
- **templateParams**: Template params to pass when processing the template manifest. Each param name should be followed by its value. Eg `[“NAME”,“my-vm”,“DESC”,“blue”]`
- **dataVolumes**: Add DVs to VM Volumes.
- **ownDataVolumes**: Add DVs to VM Volumes and add VM to DV ownerReferences. These DataVolumes will be deleted once the created VM gets deleted.
- **persistentVolumeClaims**: Add PVCs to VM Volumes.
- **ownPersistentVolumeClaims**: Add PVCs to VM Volumes and add VM to PVC ownerReferences. These PVCs will be deleted once the created VM gets deleted.

### Results

- **name**: The name of VM that was created.
- **namespace**: The namespace of VM that was created.

### Usage

Please see [examples](examples)
