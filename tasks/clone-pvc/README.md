# Clone PVC Task

This task clones PersistentVolumeClaims to target PVC by using CDI DataVolumes.

### Service Account

This task should be run with `clone-pvc-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **sourcePVC**: Source PersistentVolumeClaim to clone from.
- **sourcePVCNamespace**: Namespace of sourcePVC. (defaults to active namespace)
- **targetPVC**: Target PersistentVolumeClaim to clone to. The name will be generated if not specified. Implemented by CDI DataVolumes.
- **targetPVCNamespace**: Namespace of targetPVC. (defaults to active namespace)
- **storageClass**: StorageClass to use when cloning. (defaults to empty; ie default StorageClass)
- **waitForSuccess**: Set to `true` or `false` if container should wait for Ready condition of a DataVolume.
  
### Results

- **name**: The name of target PVC.
- **namespace**: The namespace of a target PVC.

### Usage

Please see [examples](examples)
