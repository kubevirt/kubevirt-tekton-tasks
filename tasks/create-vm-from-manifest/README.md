# Create VirtualMachine from Manifest Task

This task creates a VirtualMachine from YAML manifest

### Service Account

This task should be run with `create-vm-from-manifest-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **manifest**: YAML manifest of a VirtualMachine resource to be created.
- **namespace**: Namespace where to create the VM. (defaults to manifest namespace or active namespace)
- **startVM**: Start vm after creation.
- **dataVolumes**: Add DVs to VM Volumes. Replaces a particular volume if in VOLUME_NAME:DV_NAME format. Eg. `["rootdisk:my-dv", "my-dv2"]`
- **ownDataVolumes**: Add DVs to VM Volumes and add VM to DV ownerReferences. These DataVolumes will be deleted once the created VM gets deleted. Replaces a particular volume if in VOLUME_NAME:DV_NAME format. Eg. `["rootdisk:my-dv", "my-dv2"]`
- **persistentVolumeClaims**: Add PVCs to VM Volumes. Replaces a particular volume if in VOLUME_NAME:PVC_NAME format. Eg. `["rootdisk:my-pvc", "my-pvc2"]`
- **ownPersistentVolumeClaims**: Add PVCs to VM Volumes and add VM to PVC ownerReferences. These PVCs will be deleted once the created VM gets deleted. Replaces a particular volume if in VOLUME_NAME:PVC_NAME format. Eg. `["rootdisk:my-pvc", "my-pvc2"]`

### Results

- **name**: The name of a VM that was created.
- **namespace**: The namespace of a VM that was created.

### Usage

Please see [examples](examples) on how to create VMs.
