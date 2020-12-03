# Clone OpenShift Template Task

This task clones an OpenShift Template.

### Service Account

This task should be run with `clone-template-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **sourceTemplate**: Source Template to clone from.
- **sourceTemplateNamespace**: Namespace of sourceTemplate. (defaults to active namespace)
- **targetTemplate**: Target Template to clone to. The name will be generated if not specified.
- **targetTemplateNamespace**: Namespace of targetTemplate. (defaults to active namespace)
- **clonePVCs**: Clone all PVCs of the sourceTemplate.
  
### Results

- **name**: The name of target Template.
- **namespace**: The namespace of a target Template.

### Usage

Please see [examples](examples)
