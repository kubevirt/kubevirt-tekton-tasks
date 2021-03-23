# Disk Virt Sysprep Task

This task uses [virt-sysprep](https://libguestfs.org/virt-sysprep.1.html) to run a sysprep script on a target pvc.

{% if task_yaml.metadata.annotations['task.kubevirt.io/associatedServiceAccount'] is defined %}
### Service Account

This task should be run with `{{task_yaml.metadata.annotations['task.kubevirt.io/associatedServiceAccount']}}` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.
{% endif %}

### Parameters

{% for item in task_yaml.spec.params %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}


### Usage

Please see [examples](examples)

#### Common Errors

- The input PVC disk should not be accessed by a running VM or other tools like virt-sysprep task concurrently.
The task will fail with a generic `...guestfs_launch failed...` message.
A verbose parameter can be set to true for more information.

### OS support

- Linux: full; all the sysprep commands work
- Windows: partial; only some sysprep commands work
