{% if is_cleanup %}
# Cleanup VM Task

This task can execute a script, or a command in a Virtual Machine and stop/delete 
the VM afterwards. Best used together with tekton pipelines finally construct.
{% else %}
# Execute in VM Task

This task can execute a script, or a command in a Virtual Machine
{% endif %}

## `{{ task_name }}`

### Installation

Install the Task

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}.yaml
```

Install one of the following rbac permissions to the active namespace
{% if is_cleanup %}
  - Permissions for executing/stopping/deleting VMs from active namespace
{% else %}
  - Permissions for executing in VMs from active namespace
{% endif %}
    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}-namespace-rbac.yaml
    ```
{% if is_cleanup %}
  - Permissions for executing/stopping/deleting VMs from the cluster
{% else %}
  - Permissions for executing in VMs from the cluster
{% endif %}
    ```bash
    TARGET_NAMESPACE="$(kubectl config view --minify --output 'jsonpath={..namespace}')"
    wget -qO - https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/{{ task_name }}/manifests/{{ task_name }}-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
    ```

### Service Account

This task should be run with `{{main_task_yaml.metadata.annotations['task.kubevirt.io/associatedServiceAccount']}}` serviceAccount.

### Parameters

{% for item in main_task_yaml.spec.params %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}

### Secret format

The secret is used for storing credentials and options used in VM authentication. Following fields are recognized: 

- **type** (required): One of: ssh.
##### SSH section
- **user**: User to log in as.
- **private-key**: Private key to use for authentication.
- **host-public-key**: Public key of known host to connect to.
- **disable-strict-host-key-checking**: host-public-key (authorized-key) does not have to be supplied when this value is set to true.
- **additional-ssh-options**: Additional arguments to pass to the SSH command.

Please see [secret](examples/secrets) examples 

### Usage

Please see [examples](examples)
