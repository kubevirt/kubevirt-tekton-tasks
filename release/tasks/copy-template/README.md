# Copy okd template Task

#### Task is deprecated and will be removed in future versions.

This task copies a KubeVirt Virtual Machine template. 
A bundle of predefined templates to use can be found in [Common Templates](https://github.com/kubevirt/common-templates) project.

### Parameters

- **sourceTemplateName**: Name of an OpenShift template to copy template from.
- **sourceTemplateNamespace**: Namespace of an source OpenShift template to copy template from. (defaults to active namespace)
- **targetTemplateName**: Name of an target OpenShift template.
- **targetTemplateNamespace**: Namespace of an target OpenShift template to create in. (defaults to active namespace)
- **allowReplace**: Allow replacing already existing template (same combination name/namespace). Allowed values true/false
- **setOwnerReference**: Set owner reference to the new object created by the task run pod. Allowed values true/false

### Results

- **name**: The name of a template that was created.
- **namespace**: The namespace of a template that was created.

### Usage

Task run using resolver:
```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
    generateName: copy-template-taskrun-resolver-
spec:
    params:
    -   name: sourceTemplateName
        value: source-vm-template-example
    -   name: targetTemplateName
        value: target-vm-template-example
    taskRef:
        params:
        -   name: catalog
            value: kubevirt-tekton-tasks
        -   name: type
            value: artifact
        -   name: kind
            value: task
        -   name: name
            value: copy-template
        -   name: version
            value: v0.22.0
        resolver: hub
```

### Usage in different namespaces

You can use task to do actions in different namespace. To do that, tasks requires special permissions. Apply these RBAC objects and permissions and update accordingly task run object with correct serviceAccount:

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: copy-template-task
rules:
-   apiGroups:
    - template.openshift.io
    resources:
    - templates
    verbs:
    - get
    - list
    - watch
    - create
    - update
---
apiVersion: v1
kind: ServiceAccount
metadata:
    name: copy-template-task
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
    name: copy-template-task
roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: copy-template-task
subjects:
-   kind: ServiceAccount
    name: copy-template-task
---
```

### Platforms

The Task can be run on linux/amd64 platform.
