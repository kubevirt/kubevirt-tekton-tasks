# Create VirtualMachine from OpenShift Template Task

This task creates a VirtualMachine from OpenShift Template.
Virtual machines can be described and parametrized in a generic form with these templates.
A bundle of predefined templates to use can be found in [Common Templates](https://github.com/kubevirt/common-templates) project.

### Installation

Install the `create-vm-from-template` task in active namespace. You have to repeat this for every namespace in which you'd like to run the tasks.

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-vm-from-template/manifests/create-vm-from-template.yaml
```

For more information on how to utilize this task in different namespaces, please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md).

### Service Account

This task should be run with `create-vm-from-template-task` serviceAccount.

### Parameters

- **templateName**: Name of an OpenShift template to create VM from.
- **templateNamespace**: Namespace of an OpenShift template to create VM from. (defaults to active namespace)
- **templateParams**: Template params to pass when processing the template manifest. Each param should have KEY:VAL format. Eg `["NAME:my-vm", "DESC:blue"]`
- **vmNamespace**: Namespace where to create the VM. (defaults to active namespace)
- **dataVolumes**: Add DVs to VM Volumes.
- **ownDataVolumes**: Add DVs to VM Volumes and add VM to DV ownerReferences. These DataVolumes will be deleted once the created VM gets deleted.
- **persistentVolumeClaims**: Add PVCs to VM Volumes.
- **ownPersistentVolumeClaims**: Add PVCs to VM Volumes and add VM to PVC ownerReferences. These PVCs will be deleted once the created VM gets deleted.

### Results

- **name**: The name of a VM that was created.
- **namespace**: The namespace of a VM that was created.

### Usage

Please see [examples](examples) on how to create VMs from a template.
