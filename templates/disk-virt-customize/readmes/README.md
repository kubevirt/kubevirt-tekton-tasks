# Virt Customize Task

This task uses `virt-customize` to run a customize script on a target pvc.

{% if task_yaml.metadata.annotations['task.kubevirt.io/associatedServiceAccount'] is defined %}
### Service Account

This task should be run with `{{task_yaml.metadata.annotations['task.kubevirt.io/associatedServiceAccount']}}` serviceAccount.
{% endif %}

### Parameters

{% for item in task_yaml.spec.params %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}


### Usage

Please see [examples](examples)
