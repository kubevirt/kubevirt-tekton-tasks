# Wait For a VirtualMachineInstance Status Task

This task waits for a specific status of a VirtualMachineInstance (VMI) and fails/succeeds accordingly.

### Service Account

This task should be run with `wait-for-vmi-status-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **vmiName**: Name of a VirtualMachineInstance to wait for.
- **vmiNamespace**: Namespace of a VirtualMachineInstance to wait for. (defaults to manifest namespace or active namespace)
- **successCondition**: A label selector expression to decide if the VirtualMachineInstance (VMI) is in a success state. Eg. `status.phase == Succeeded`. JSONPath format of a parametre is also supported and can be used like this `jsonpath='{.status.phase}' == Succeeded`. It is evaluated on each VMI update and will result in this task succeeding if true. It uses kubernetes label selection syntax and can be applied against any field of the resource (not just labels). Multiple AND conditions can be represented by comma delimited expressions. For more details, see: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/.
- **failureCondition**: A label selector expression to decide if the VirtualMachineInstance (VMI) is in a failed state. Eg. `status.phase in (Failed, Unknown)`. JSONPath format of a parametre is also supported and can be used like this `jsonpath='{.status.phase}' == Failed`. It is evaluated on each VMI update and will result in this task failing if true. It uses kubernetes label selection syntax and can be applied against any field of the resource (not just labels). Multiple AND conditions can be represented by comma delimited expressions. For more details, see: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/.

### Usage

Please see [examples](examples)
