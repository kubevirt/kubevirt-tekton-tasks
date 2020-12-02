# Create DataVolume Task

This task creates a DataVolume with oc client.

## `create-datavolume-from-manifest`

### Installation

Install the Task

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-datavolume/manifests/create-datavolume-from-manifest.yaml
```

Install one of the following rbac permissions to the active namespace
- Permissions for creating DataVolumes in active namespace
  ```bash
  kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-datavolume/manifests/create-datavolume-namespace-rbac.yaml
  ```
- Permissions for creating DataVolumes in the cluster
  ```bash
  TARGET_NAMESPACE="$(kubectl config current-context | cut -d/ -f1)"
  wget -qO - https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-datavolume/manifests/create-datavolume-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
  ```

### Service Account

This task should be run with `create-datavolume-task` serviceAccount.

### Parameters

- **manifest**: YAML manifest of a DataVolume resource to be created.
- **waitForSuccess**: Set to `true` or `false` if container should wait for Ready condition of a DataVolume.
  
### Results

- **name**: The name of DataVolume that was created.
- **namespace**: The namespace of DataVolume that was created.

### Usage

Please see [examples](examples)