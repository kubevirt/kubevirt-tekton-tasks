# Create VirtualMachine from Manifest Task

This task creates a VirtualMachine from YAML manifest

### Parameters

- **manifest**: YAML manifest of a VirtualMachine resource to be created.
- **virtctl**: Parameters for virtctl create vm command that will be used to create VirtualMachine.
- **namespace**: Namespace where to create the VM. (defaults to manifest namespace or active namespace)
- **startVM**: Set to true or false to start / not start vm after creation. In case of runStrategy is set to Always, startVM flag is ignored.
- **runStrategy**: Set runStrategy to VM. If runStrategy is set, vm.spec.running attribute is set to nil.
- **setOwnerReference**: Set owner reference to the new object created by the task run pod. Allowed values true/false

### Results

- **name**: The name of a VM that was created.
- **namespace**: The namespace of a VM that was created.

### Usage

Task run using resolver:
```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
    generateName: create-vm-from-manifest-taskrun-resolver-
spec:
    params:
    -   name: manifest
        value: <VM manifest>
    taskRef:
        params:
        -   name: catalog
            value: kubevirt-tekton-tasks
        -   name: type
            value: artifact
        -   name: kind
            value: task
        -   name: name
            value: create-vm-from-manifest
        -   name: version
            value: v0.22.0
        resolver: hub
```

As an example for `manifest` parameter, you can use this VM definition:
```
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
    generateName: vm-fedora-
    labels:
        kubevirt.io/vm: vm-fedora
spec:
    running: false
    template:
        metadata:
            labels:
                kubevirt.io/vm: vm-fedora
        spec:
            domain:
                devices:
                    disks:
                    -   disk:
                            bus: virtio
                        name: containerdisk
                    -   disk:
                            bus: virtio
                        name: cloudinitdisk
                memory:
                    guest: 1Gi
            volumes:
            -   containerDisk:
                    image: quay.io/containerdisks/fedora:latest
                name: containerdisk
            -   cloudInitNoCloud:
                    userData: '#!/bin/sh

                        echo ''printed from cloud-init userdata''

                        '
                name: cloudinitdisk
```

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

### Platforms

The Task can be run on linux/amd64 platform.
