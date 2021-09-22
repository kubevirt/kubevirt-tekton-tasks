# Modify OpenShift template Task

This task modifies template.
A bundle of predefined templates to use can be found in [Common Templates](https://github.com/kubevirt/common-templates) project.

### Service Account

This task should be run with `modify-vm-template-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **templateName**: Name of an OpenShift template.
- **templateNamespace**: Namespace of an source OpenShift template. (defaults to active namespace)
- **cpuSockets**: Number of CPU sockets
- **cpuCores**: Number of CPU cores
- **cpuThreads**: Number of CPU threads
- **memory**: Number of memory vm can use
- **templateLabels**: Template labels. If template contains same label, it will be replaced.
- **templateAnnotations**: Template Annotations. If template contains same annotation, it will be replaced.
- **vmLabels**: VM labels. If VM contains same label, it will be replaced.
- **vmAnnotations**: VM annotations. If VM contains same annotation, it will be replaced.

### Results

- **name**: The name of a template that was updated.
- **namespace**: The namespace of a template that was updated.

### Usage

Please see [examples](examples) on how to do a copy template from a template.
