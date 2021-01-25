# Create DataVolume from Manifest Task

This task creates a DataVolume with oc client.

### Installation

Install the `create-datavolume-from-manifest` task

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-datavolume-from-manifest/manifests/create-datavolume-from-manifest.yaml
```

Install one of the following rbac permissions to the active namespace
- Permissions for creating DataVolumes in active namespace
  ```bash
  kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-datavolume-from-manifest/manifests/create-datavolume-from-manifest-namespace-rbac.yaml
  ```
- Permissions for creating DataVolumes in the cluster
  ```bash
  TARGET_NAMESPACE="$(kubectl config view --minify --output 'jsonpath={..namespace}')"
  wget -qO - https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/create-datavolume-from-manifest/manifests/create-datavolume-from-manifest-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
  ```

### Service Account

This task should be run with `create-datavolume-from-manifest-task` serviceAccount.

### Parameters

- **manifest**: YAML manifest of a DataVolume resource to be created.
- **waitForSuccess**: Set to `true` or `false` if container should wait for Ready condition of a DataVolume.
  
### Results

- **name**: The name of DataVolume that was created.
- **namespace**: The namespace of DataVolume that was created.

### Usage

Please see [examples](examples)
