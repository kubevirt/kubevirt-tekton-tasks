# Create VirtualMachine from Manifest Task

This task creates a VirtualMachine from YAML manifest

### Installation

Install the `{{ task_name }}` task

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}.yaml
```

Install one of the following rbac permissions to the active namespace
- Permissions for using VMs and storage in active namespace
  ```bash
  kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}-namespace-rbac.yaml
  ```
- Permissions for using VMs and storage in the cluster
  ```bash
  TARGET_NAMESPACE="$(kubectl config view --minify --output 'jsonpath={..namespace}')"
  wget -qO - https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
  ```

### Service Account

This task should be run with `{{task_yaml.metadata.annotations['task.kubevirt.io/associatedServiceAccount']}}` serviceAccount.

### Parameters

{% for item in task_yaml.spec.params %}
{% if item.name == "templateParams" %}
- **{{ item.name }}**: {{ item.description | replace('[', '`[')   | replace(']', ']`')}}
{% else %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endif %}
{% endfor %}

### Results

{% for item in task_yaml.spec.results %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}

### Usage

Please see [examples](examples) on how to create VMs.
