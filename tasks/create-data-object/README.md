# Create Data Object Task

This task creates a data object (DataVolumes or DataSources).

### Service Account

This task should be run with `create-data-object-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **manifest**: YAML manifest of a data object to be created.
- **namespace**: Namespace where to create the data object. (defaults to manifest namespace or active namespace)
- **waitForSuccess**: Set to `true` or `false` if container should wait for Ready condition of the data object.
- **allowReplace**: Allow replacing an already existing data object (same combination name/namespace). Allowed values true/false
  
### Results

- **name**: The name of the data object that was created.
- **namespace**: The namespace of the data object that was created.

### Usage

Please see [examples](examples) on how to create data objects.
