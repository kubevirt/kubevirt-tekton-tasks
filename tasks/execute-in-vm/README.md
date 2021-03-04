# Execute in VM Task

This task can execute a script, or a command in a Virtual Machine

### Installation

Install the `execute-in-vm` task

```bash
kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/execute-in-vm/manifests/execute-in-vm.yaml
```

Install one of the following rbac permissions to the active namespace
  - Permissions for executing in VMs from active namespace
    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/execute-in-vm/manifests/execute-in-vm-namespace-rbac.yaml
    ```
  - Permissions for executing in VMs from the cluster
    ```bash
    TARGET_NAMESPACE="$(kubectl config view --minify --output 'jsonpath={..namespace}')"
    wget -qO - https://raw.githubusercontent.com/kubevirt/kubevirt-tekton-tasks/main/tasks/execute-in-vm/manifests/execute-in-vm-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
    ```

### Service Account

This task should be run with `execute-in-vm-task` serviceAccount.

### Parameters

- **vmName**: Name of a VM to execute the action in.
- **vmNamespace**: Namespace of a VM to execute the action in. (defaults to active namespace)
- **secretName**: Secret to use when connecting to a VM.
- **command**: Command to execute in a VM.
- **args**: Arguments of a command.
- **script**: Script to execute in a VM.

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

- [start postgresql service over ssh](examples/taskruns/execute-in-vm-with-ssh-taskrun.yaml)
