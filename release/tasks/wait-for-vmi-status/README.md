# Wait For a VirtualMachineInstance Status Task

This task waits for a specific status of a VirtualMachineInstance (VMI) and fails/succeeds accordingly.

### Parameters

- **vmiName**: Name of a VirtualMachineInstance to wait for.
- **vmiNamespace**: Namespace of a VirtualMachineInstance to wait for. (defaults to manifest namespace or active namespace)
- **successCondition**: A label selector expression to decide if the VirtualMachineInstance (VMI) is in a success state. Eg. `status.phase == Succeeded`. It is evaluated on each VMI update and will result in this task succeeding if true. It uses kubernetes label selection syntax and can be applied against any field of the resource (not just labels). Multiple AND conditions can be represented by comma delimited expressions. For more details, see: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/.
- **failureCondition**: A label selector expression to decide if the VirtualMachineInstance (VMI) is in a failed state. Eg. `status.phase in (Failed, Unknown)`. It is evaluated on each VMI update and will result in this task failing if true. It uses kubernetes label selection syntax and can be applied against any field of the resource (not just labels). Multiple AND conditions can be represented by comma delimited expressions. For more details, see: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/.

### Usage

Task run using resolver:
```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
    generateName: wait-for-vmi-status-taskrun-resolver-
spec:
    params:
    -   name: vmiName
        value: example-vm
    -   name: successCondition
        value: status.phase == Succeeded
    -   name: failureCondition
        value: status.phase in (Failed, Unknown)
    taskRef:
        params:
        -   name: catalog
            value: kubevirt-tekton-tasks
        -   name: type
            value: artifact
        -   name: kind
            value: task
        -   name: name
            value: wait-for-vmi-status
        -   name: version
            value: v0.21.0
        resolver: hub
```

### Usage in different namespaces

You can use task to do actions in different namespace. To do that, tasks requires special permissions. Apply these RBAC objects and permissions and update accordingly task run object with correct serviceAccount:

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: wait-for-vmi-status-task
rules:
-   apiGroups:
    - kubevirt.io
    resources:
    - virtualmachineinstances
    verbs:
    - get
    - list
    - watch
---
apiVersion: v1
kind: ServiceAccount
metadata:
    name: wait-for-vmi-status-task
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
    name: wait-for-vmi-status-task
roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: wait-for-vmi-status-task
subjects:
-   kind: ServiceAccount
    name: wait-for-vmi-status-task
---
```

### Platforms

The Task can be run on linux/amd64 platform.
