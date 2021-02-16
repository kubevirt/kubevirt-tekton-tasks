# Create VirtualMachine from Manifest Task

This task creates a VirtualMachine from YAML manifest

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
