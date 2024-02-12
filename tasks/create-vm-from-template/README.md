# Create VirtualMachine from OKD Template Task

#### Task is deprecated and will be removed in future versions.

This task creates a VirtualMachine from OKD Template.
Virtual machines can be described and parametrized in a generic form with these templates.
A bundle of predefined templates to use can be found in [Common Templates](https://github.com/kubevirt/common-templates) project.

### Parameters

- **templateName**: Name of an OKD template to create VM from.
- **templateNamespace**: Namespace of an OKD template to create VM from. (defaults to active namespace)
- **templateParams**: Template params to pass when processing the template manifest. Each param should have KEY:VAL format. Eg `["NAME:my-vm", "DESC:blue"]`
- **vmNamespace**: Namespace where to create the VM. (defaults to active namespace)
- **startVM**: Set to true or false to start / not start vm after creation. In case of runStrategy is set to Always, startVM flag is ignored.
- **runStrategy**: Set runStrategy to VM. If runStrategy is set, vm.spec.running attribute is set to nil.

### Results

- **name**: The name of a VM that was created.
- **namespace**: The namespace of a VM that was created.

### Usage

Please see [examples](examples) on how to create VMs from a template.

### Usage in different namespaces

You can use task to do actions in different namespace. To do that, tasks requires special permissions. Apply these RBAC objects and permissions and update accordingly task run object with correct serviceAccount:

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: create-vm-from-template-task
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
    - kubevirt.io
    resources:
    - virtualmachines/finalizers
    verbs:
    - get
-   apiGroups:
    - template.openshift.io
    resources:
    - templates
    verbs:
    - get
    - list
    - watch
-   apiGroups:
    - template.openshift.io
    resources:
    - processedtemplates
    verbs:
    - create
-   apiGroups:
    - cdi.kubevirt.io
    resources:
    - datavolumes
    verbs:
    - create
-   apiGroups:
    - subresources.kubevirt.io
    resources:
    - virtualmachines/start
    verbs:
    - update
---
apiVersion: v1
kind: ServiceAccount
metadata:
    name: create-vm-from-template-task
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
    name: create-vm-from-template-task
roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: create-vm-from-template-task
subjects:
-   kind: ServiceAccount
    name: create-vm-from-template-task
---
```
