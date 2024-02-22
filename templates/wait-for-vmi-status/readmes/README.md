# Wait For a VirtualMachineInstance Status Task

This task waits for a specific status of a VirtualMachineInstance (VMI) and fails/succeeds accordingly.

### Parameters

{% for item in task_yaml.spec.params %}
{% if 'Condition' in item.name %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }} It uses kubernetes label selection syntax and can be applied against any field of the resource (not just labels). Multiple AND conditions can be represented by comma delimited expressions. For more details, see: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/.
{% else %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endif %}
{% endfor %}

### Usage

Task run using resolver:
```
{{ task_run_resolver_yaml | to_nice_yaml }}```

### Usage in different namespaces

You can use task to do actions in different namespace. To do that, tasks requires special permissions. Apply these RBAC objects and permissions and update accordingly task run object with correct serviceAccount:

```
{% for item in rbac_yaml %}
{{ item | to_nice_yaml }}---
{% endfor %}
```

### Platforms

The Task can be run on linux/amd64 platform.
