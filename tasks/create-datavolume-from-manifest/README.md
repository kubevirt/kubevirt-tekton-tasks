# Create DataVolume from Manifest Task

This task creates a DataVolume with oc client.

### Installation

Install the `create-datavolume-from-manifest` task in active namespace. You have to repeat this for every namespace in which you'd like to run the tasks.

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-datavolume-from-manifest/manifests/create-datavolume-from-manifest.yaml
```

For more information on how to utilize this task in different namespaces, please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md).

### Service Account

This task should be run with `create-datavolume-from-manifest-task` serviceAccount.

### Parameters

- **manifest**: YAML manifest of a DataVolume resource to be created.
- **waitForSuccess**: Set to `true` or `false` if container should wait for Ready condition of a DataVolume.
  
### Results

- **name**: The name of DataVolume that was created.
- **namespace**: The namespace of DataVolume that was created.

### Usage

Please see [examples](examples) on how to create DataVolumes.
