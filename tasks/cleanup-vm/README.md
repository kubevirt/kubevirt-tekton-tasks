# Cleanup VM Task

This task can execute a script, or a command in a Virtual Machine and stop/delete 
the VM afterwards. Best used together with tekton pipelines finally construct.

## `cleanup-vm`

### Installation

Install the Task

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/cleanup-vm/manifests/cleanup-vm.yaml
```

Install one of the following rbac permissions to the active namespace
  - Permissions for executing/stopping/deleting VMs from active namespace
    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/cleanup-vm/manifests/cleanup-vm-namespace-rbac.yaml
    ```
  - Permissions for executing/stopping/deleting VMs from the cluster
    ```bash
    TARGET_NAMESPACE="$(kubectl config current-context | cut -d/ -f1)"
    wget -qO - https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/cleanup-vm/manifests/cleanup-vm-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
    ```

### Service Account

This task should be run with `cleanup-vm-task` serviceAccount.

### Parameters

- **vmName**: Name of a VM to execute the action in.
- **vmNamespace**: Namespace of a VM to execute the action in. (defaults to active namespace)
- **stop**: Stops the VM after executing the commands when set to true.
- **delete**: Deletes the VM after executing the commands when set to true.
- **timeout**: Timeout for the command/script (includes potential VM start). The VM will be stoped or deleted accordingly once the timout expires. Should be in a 3h2m1s format.
- **secretName**: Secret to use when connecting to a VM.
- **command**: Command to execute in a VM.
- **args**: Arguments of a command.
- **script**: Script to execute in a VM.

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
