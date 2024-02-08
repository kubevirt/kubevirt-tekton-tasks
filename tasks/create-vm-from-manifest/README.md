# Create VirtualMachine from Manifest Task

This task creates a VirtualMachine from YAML manifest

### Parameters

- **manifest**: YAML manifest of a VirtualMachine resource to be created.
- **virtctl**: Parameters for virtctl create vm command that will be used to create VirtualMachine.
- **namespace**: Namespace where to create the VM. (defaults to manifest namespace or active namespace)
- **startVM**: Set to true or false to start / not start vm after creation. In case of runStrategy is set to Always, startVM flag is ignored.
- **runStrategy**: Set runStrategy to VM. If runStrategy is set, vm.spec.running attribute is set to nil.

### Results

- **name**: The name of a VM that was created.
- **namespace**: The namespace of a VM that was created.

### Usage

Please see [examples](examples) on how to create VMs.

### Usage in different namespaces

You can use task to do actions in different namespace. To do that, tasks requires special permissions. Apply these RBAC objects and permissions and update accordingly task run object with correct serviceAccount:

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: create-vm-from-manifest-task
rules:
-   apiGroups:
    - kubevirt.io
    resources:
    - virtualmachines
    - virtualmachineinstances
    verbs:
    - get
    - list
    - watch
    - create
    - update
-   apiGroups:
    - subresources.kubevirt.io
    resources:
    - virtualmachines/start
    verbs:
    - update
-   apiGroups:
    - ''
    resources:
    - persistentvolumeclaims
    verbs:
    - update
---
apiVersion: v1
kind: ServiceAccount
metadata:
    name: create-vm-from-manifest-task
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
    name: create-vm-from-manifest-task
roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: create-vm-from-manifest-task
subjects:
-   kind: ServiceAccount
    name: create-vm-from-manifest-task
---
```
