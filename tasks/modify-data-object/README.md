# Modify Data Object Task

This task modifies a data object (DataVolumes or DataSources).

### Service Account

This task should be run with `modify-data-object-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **manifest**: YAML manifest of a data object to be created.
- **namespace**: Namespace where to create the data object. (defaults to manifest namespace or active namespace)
- **waitForSuccess**: Set to `true` or `false` if container should wait for Ready condition of the data object.
- **allowReplace**: Allow replacing an already existing data object (same combination name/namespace). Allowed values true/false
- **deleteObject**: Set to `true` or `false` if task should delete the specified datavolume or datasource. If set to 'true' the ds/dv will be deleted and all other parameters are ignored.
- **deleteObjectKind**: Kind of the data object to delete. This parameter is used only for Delete operation.
- **deleteObjectName**: Name of the data object to delete. This parameter is used only for Delete operation.
  
### Results

- **name**: The name of the data object that was created.
- **namespace**: The namespace of the data object that was created.

### Usage

Please see [examples](examples) on how to modify data objects.
