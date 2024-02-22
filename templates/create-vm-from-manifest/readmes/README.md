# Create VirtualMachine from Manifest Task

This task creates a VirtualMachine from YAML manifest

### Parameters

{% for item in task_yaml.spec.params %}
{% if 'Volume' in item.name %}
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

Task run using resolver:
```
{{ task_run_resolver_yaml | to_nice_yaml }}```

As an example for `manifest` parameter, you can use this VM definition:
```
{{ vm_manifest | to_nice_yaml }}```

### Usage in different namespaces

You can use task to do actions in different namespace. To do that, tasks requires special permissions. Apply these RBAC objects and permissions and update accordingly task run object with correct serviceAccount:

```
{% for item in rbac_yaml %}
{{ item | to_nice_yaml }}---
{% endfor %}
```

### Platforms

The Task can be run on linux/amd64 platform.
