# Create DataVolume Task

This task creates a DataVolume with oc client.

## `{{ create_dv_from_manifest_name }}`

### Installation

Install the Task

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ create_dv_from_manifest_name }}.yaml
```

Install one of the following rbac permissions to the active namespace
- Permissions for creating DataVolumes in active namespace
  ```bash
  kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}-namespace-rbac.yaml
  ```
- Permissions for creating DataVolumes in the cluster
  ```bash
  TARGET_NAMESPACE="$(kubectl config current-context | cut -d/ -f1)"
  wget -qO - https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
  ```

### Service Account

This task should be run with `{{create_dv_from_manifest_yaml.metadata.annotations['task.kubevirt.io/associatedServiceAccount']}}` serviceAccount.

### Parameters

{% for item in create_dv_from_manifest_yaml.spec.params %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}
  
### Results

{% for item in create_dv_from_manifest_yaml.spec.results %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}

### Usage

Please see [examples](examples)