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

#### Common Errors

- The input PVC disk should not be accessed by a running VM or other tools like virt-customize task concurrently.
The task will fail with a generic `...guestfs_launch failed...` message.
Verbose parameter can be set to true for more information.

### OS support

- Linux: full; all the customize commands work
- Windows: partial; only some customize commands work
