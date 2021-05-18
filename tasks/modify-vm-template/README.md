# Modify Virtual Machine OpenShift Template Task

This task modifies a VM inside a template and template metadata.

### Service Account

This task should be run with `modify-vm-template-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **template**: Template to modify.
- **templateNamespace**: Namespace of template. (defaults to active namespace)
- **cpus**: Changes number of VM CPUs.
- **memory**: Changes memory of the VM.
- **providerName**: Specifies the name of the template provider.
- **supportURL**: Specifies the support URL of the template.

### Usage

Please see [examples](examples)
