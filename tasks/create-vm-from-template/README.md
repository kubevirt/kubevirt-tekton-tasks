# Create VirtualMachine from OKD Template Task

This task creates a VirtualMachine from OKD Template.
Virtual machines can be described and parametrized in a generic form with these templates.
A bundle of predefined templates to use can be found in [Common Templates](https://github.com/kubevirt/common-templates) project.

### Service Account

This task should be run with `create-vm-from-template-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **templateName**: Name of an OKD template to create VM from.
- **templateNamespace**: Namespace of an OKD template to create VM from. (defaults to active namespace)
- **templateParams**: Template params to pass when processing the template manifest. Each param should have KEY:VAL format. Eg `["NAME:my-vm", "DESC:blue"]`
- **vmNamespace**: Namespace where to create the VM. (defaults to active namespace)
- **startVM**: Set to true or false to start / not start vm after creation. In case of runStrategy is set to Always, startVM flag is ignored.
- **runStrategy**: Set runStrategy to VM. If runStrategy is set, vm.spec.running attribute is set to nil.
- **dataVolumes**: Add DVs to VM Volumes. Replaces a particular volume if in VOLUME_NAME:DV_NAME format. Eg. `["rootdisk:my-dv", "my-dv2"]`
- **ownDataVolumes**: Add DVs to VM Volumes and add VM to DV ownerReferences. These DataVolumes will be deleted once the created VM gets deleted. Replaces a particular volume if in VOLUME_NAME:DV_NAME format. Eg. `["rootdisk:my-dv", "my-dv2"]`
- **persistentVolumeClaims**: Add PVCs to VM Volumes. Replaces a particular volume if in VOLUME_NAME:PVC_NAME format. Eg. `["rootdisk:my-pvc", "my-pvc2"]`
- **ownPersistentVolumeClaims**: Add PVCs to VM Volumes and add VM to PVC ownerReferences. These PVCs will be deleted once the created VM gets deleted. Replaces a particular volume if in VOLUME_NAME:PVC_NAME format. Eg. `["rootdisk:my-pvc", "my-pvc2"]`

### Results

- **name**: The name of a VM that was created.
- **namespace**: The namespace of a VM that was created.

### Usage

Please see [examples](examples) on how to create VMs from a template.
