{% if is_cleanup %}
# Cleanup VM Task

This task can execute a script, or a command in a Virtual Machine and stop/delete 
the VM afterwards. Best used together with tekton pipelines finally construct.
{% else %}
# Execute in VM Task

This task can execute a script, or a command in a Virtual Machine
{% endif %}

### Installation

Install the `{{ task_name }}` task in active namespace. You have to repeat this for every namespace in which you'd like to run the tasks.

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}.yaml
```

For more information on how to utilize this task in different namespaces, please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md).

### Service Account

This task should be run with `{{task_yaml.metadata.annotations['task.kubevirt.io/associatedServiceAccount']}}` serviceAccount.

### Parameters

{% for item in task_yaml.spec.params %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}

### Secret format

The secret is used for storing credentials and options used in VM authentication.

##### Specifying a type

The secret should be of one of the following types:

- `kubernetes.io/ssh-auth`
- `Opaque`: Secret data should include the following key.
    - **type**: One of: ssh.

##### SSH section

Following secret data keys are recognized for SSH connections:

- **user**: User to log in as.
- **ssh-privatekey**: Private key to use for authentication.
- **host-public-key**: Public key of known host to connect to.
- **disable-strict-host-key-checking**: host-public-key (authorized-key) does not have to be supplied when this value is set to true.
- **additional-ssh-options**: Additional arguments to pass to the SSH command.

Please see [secret](examples/secrets) examples.

### Usage

Please see [examples](examples).

#### Specific examples

{% if is_cleanup %}
- [delete a VM](examples/taskruns/cleanup-vm-simple-taskrun.yaml)
- [stop postgresql service over ssh and stop a VM](examples/taskruns/cleanup-vm-with-ssh-taskrun.yaml)
{% else %}
- [start postgresql service over ssh](examples/taskruns/execute-in-vm-with-ssh-taskrun.yaml)
{% endif %}
