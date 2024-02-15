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
  
### Results

- **name**: The name of the data object that was created.
- **namespace**: The namespace of the data object that was created.

### Usage

Please see [examples](examples) on how to modify data objects.

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
