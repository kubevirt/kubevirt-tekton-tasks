# Disk Virt Customize Task

This task uses [virt-customize](https://libguestfs.org/virt-customize.1.html) to run a customize script on a target pvc.

### Parameters

{% for item in task_yaml.spec.params %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}

### Usage

Task run using resolver:
```
{{ task_run_resolver_yaml | to_nice_yaml }}```

#### Common Errors

- The input PVC disk should not be accessed by a running VM or other tools like virt-customize concurrently.
The task will fail with a generic `...guestfs_launch failed...` message.
Verbose parameter can be set to true for more information.

- The task can end with error `Permissions denied`. This error means, the disk-virt-customize cannot access the VM image due to wrong permissions set on the file. This issue can be fixed by adding `podTemplate` to the spec of the TaskRun:
```
spec:
  podTemplate:
    securityContext:
      fsGroup: 107
      runAsUser: 107
```

### OS support

- Linux: full; all the customize commands work
- Windows: partial; only some customize commands work
