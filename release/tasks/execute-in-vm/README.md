# Execute in VM Task

This task can execute a script, or a command in a Virtual Machine

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

Task run using resolver:
```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
    generateName: execute-in-vm-with-ssh-taskrun-resolver-
spec:
    params:
    -   name: vmName
        value: vm-example
    -   name: secretName
        value: ssh-secret
    -   name: command
        value:
        - systemctl
    -   name: args
        value:
        - start
        - postgresql.service
    taskRef:
        params:
        -   name: catalog
            value: kubevirt-tekton-tasks
        -   name: type
            value: artifact
        -   name: kind
            value: task
        -   name: name
            value: execute-in-vm
        -   name: version
            value: v0.21.0
        resolver: hub
```

### Usage in different namespaces

You can use task to do actions in different namespace. To do that, tasks requires special permissions. Apply these RBAC objects and permissions and update accordingly task run object with correct serviceAccount:

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: execute-in-vm-task
rules:
-   apiGroups:
    - kubevirt.io
    resources:
    - virtualmachines
    - virtualmachineinstances
    verbs:
    - get
    - list
    - watch
-   apiGroups:
    - subresources.kubevirt.io
    resources:
    - virtualmachines/start
    - virtualmachines/stop
    - virtualmachines/restart
    verbs:
    - update
---
apiVersion: v1
kind: ServiceAccount
metadata:
    name: execute-in-vm-task
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
    name: execute-in-vm-task
roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: execute-in-vm-task
subjects:
-   kind: ServiceAccount
    name: execute-in-vm-task
---
```

### Platforms

The Task can be run on linux/amd64 platform.
