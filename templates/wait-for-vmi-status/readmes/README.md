# Wait For a VirtualMachineInstance Status Task

This task waits for a specific status of a VirtualMachineInstance (VMI) and fails/succeeds accordingly.

### Service Account

This task should be run with `{{task_yaml.metadata.annotations['task.kubevirt.io/associatedServiceAccount']}}` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

{% for item in task_yaml.spec.params %}
{% if 'Condition' in item.name %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }} It uses kubernetes label selection syntax and can be applied against any field of the resource (not just labels). Multiple AND conditions can be represented by comma delimited expressions. For more details, see: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/.
{% else %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endif %}
{% endfor %}

### Usage

Please see [examples](examples)
