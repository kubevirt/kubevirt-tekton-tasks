# Create VirtualMachine Task

Flavors of this task create a VirtualMachine from different sources

## `{{ create_vm_from_template_name }}`

### Installation

Install the Task

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ create_vm_from_template_name }}.yaml
```

Install one of the following rbac permissions to the active namespace
- Permissions for using templates/VMs in active namespace
  ```bash
  kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}-namespace-rbac.yaml
  ```
- Permissions for using templates/VMs in the cluster
  ```bash
  TARGET_NAMESPACE="$(kubectl config current-context | cut -d/ -f1)"
  wget -qO - https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
  ```

### Service Account

This task should be run with `{{create_vm_from_template_yaml.metadata.annotations['task.kubevirt.io/associatedServiceAccount']}}` serviceAccount.

### Parameters

{% for item in create_vm_from_template_yaml.spec.params %}
{% if item.name == "templateParams" %}
- **{{ item.name }}**: {{ item.description | replace('[', '`[')   | replace(']', ']`')}}
{% else %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endif %}
{% endfor %}

### Results

{% for item in create_vm_from_template_yaml.spec.results %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}

### Usage

Please see [examples](examples)