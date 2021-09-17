# Copy okd template Task

This task copies a template.
A bundle of predefined templates to use can be found in [Common Templates](https://github.com/kubevirt/common-templates) project.

### Service Account

This task should be run with `copy-template-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **sourceTemplateName**: Name of an OpenShift template to copy template from.
- **sourceTemplateNamespace**: Namespace of an source OpenShift template to copy template from. (defaults to active namespace)
- **targetTemplateName**: Name of an target OpenShift template.
- **targetTemplateNamespace**: Namespace of an target OpenShift template to create in. (defaults to active namespace)

### Results

- **name**: The name of a template that was created.
- **namespace**: The namespace of a template that was created.

### Usage

Please see [examples](examples) on how to copy a template.
