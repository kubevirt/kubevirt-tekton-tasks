# Modify Data Object Task

This task modifies a data object (DataVolumes or DataSources).

### Parameters

- **manifest**: YAML manifest of a data object to be created.
- **namespace**: Namespace where to create the data object. (defaults to manifest namespace or active namespace)
- **waitForSuccess**: Set to `true` or `false` if container should wait for Ready condition of the data object.
- **allowReplace**: Allow replacing an already existing data object (same combination name/namespace). Allowed values true/false
- **deleteObject**: Set to `true` or `false` if task should delete the specified DataVolume, DataSource or PersistentVolumeClaim. If set to 'true' the ds/dv/pvc will be deleted and all other parameters are ignored.
- **deleteObjectKind**: Kind of the data object to delete. This parameter is used only for Delete operation.
- **deleteObjectName**: Name of the data object to delete. This parameter is used only for Delete operation.
- **setOwnerReference**: Set owner reference to the new object created by the task run pod. Allowed values true/false
  
### Results

- **name**: The name of the data object that was created.
- **namespace**: The namespace of the data object that was created.

### Usage

Task run using resolver:
```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
    generateName: modify-data-object-taskrun-resolver-
spec:
    params:
    -   name: waitForSuccess
        value: 'true'
    -   name: manifest
        value: <DV or DS manifest>
    taskRef:
        params:
        -   name: catalog
            value: kubevirt-tekton-tasks
        -   name: type
            value: artifact
        -   name: kind
            value: task
        -   name: name
            value: modify-data-object
        -   name: version
            value: v0.22.0
        resolver: hub
```

As an example for `manifest` parameter, you can use this DV definition:
```
apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
    generateName: example-dv-
spec:
    pvc:
        accessModes:
        - ReadWriteOnce
        resources:
            requests:
                storage: 100Mi
        volumeMode: Filesystem
    source:
        blank: {}
```

### Usage in different namespaces

You can use task to do actions in different namespace. To do that, tasks requires special permissions. Apply these RBAC objects and permissions and update accordingly task run object with correct serviceAccount:

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: modify-data-object-task
rules:
-   apiGroups:
    - cdi.kubevirt.io
    resources:
    - datavolumes
    - datasources
    verbs:
    - get
    - create
    - delete
-   apiGroups:
    - ''
    resources:
    - pods
    verbs:
    - create
-   apiGroups:
    - ''
    resources:
    - persistentvolumeclaims
    verbs:
    - get
    - delete
---
apiVersion: v1
kind: ServiceAccount
metadata:
    name: modify-data-object-task
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
    name: modify-data-object-task
roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: modify-data-object-task
subjects:
-   kind: ServiceAccount
    name: modify-data-object-task
---
```

### Platforms

The Task can be run on linux/amd64 platform.
