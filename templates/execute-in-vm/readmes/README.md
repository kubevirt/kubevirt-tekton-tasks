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
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/master/tasks/{{ task_name }}/manifests/{{ task_name }}.yaml
```

Install one of the following rbac permissions to the active namespace
{% if is_cleanup %}
  - Permissions for executing/stopping/deleting VMs from active namespace
{% else %}
  - Permissions for executing in VMs from active namespace
{% endif %}
    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/master/tasks/{{ task_name }}/manifests/{{ task_name }}-namespace-rbac.yaml
    ```
{% if is_cleanup %}
  - Permissions for executing/stopping/deleting VMs from the cluster
{% else %}
  - Permissions for executing in VMs from the cluster
{% endif %}
    ```bash
    TARGET_NAMESPACE="$(kubectl config current-context | cut -d/ -f1)"
    wget -qO - https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/master/tasks/{{ task_name }}/manifests/{{ task_name }}-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
    ```

### Parameters

- **vmName**: Name of a VM to execute the action in.
- **vmNamespace**: Namespace of a VM to execute the action in (defaults to active namespace).
{% if is_cleanup %}
- **stop**: Stops the VM after executing the commands when set to true
- **delete**: Deletes the VM after executing the commands when set to true
- **timeout**: Timeout for the command/script. The VM will be stoped or deleted accordingly once the timout expires. Should be in a 3h2m1s format.
{% endif %}
- **secretName**: Secret to use when connecting to a VM.
- **command**: Command to execute in a VM.
- **args**: Arguments of a command.
- **script**: Script to execute in a VM

### Secret format

The secret is used for storing credentials and options used in VM authentication. Following fields are recognized: 

- **type** (required): One of: ssh.
##### SSH section
- **user** (required): Username.
- **private-key**: Private key to use for authentication. Alternatively generate-client-keys can be used instead.
- **host-public-key**: Public key of known host to connect to. Alternatively generate-host-keys can be used instead.
- **disable-strict-host-key-checking**: authorized-key does not have to be supplied when this value is set to true.
- **additional-ssh-options**: Additional arguments to pass to the SSH command.
<!---
- **generate-client-keys (TBD)**: Generates authentication keys for execute-in-vm task client if set to true. Then it will try to supply the public key to VM's authorized keys.
- **generate-host-keys (TBD)**: Generates authentication keys for the VM if set to true. Then it will try to supply the private and public key to the VM.
-->

Please see [secret](examples/secrets) examples 

### Usage

Please see [examples](examples) for examples  
