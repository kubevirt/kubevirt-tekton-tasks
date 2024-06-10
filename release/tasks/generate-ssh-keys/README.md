# Generate SSH Keys Task 

This task uses `ssh-keygen` to generate a private and public key pair

### Parameters

- **publicKeySecretName**: Name of a new or existing secret to append the generated public key to. The name will be generated and new secret created if not specified.
- **publicKeySecretNamespace**: Namespace of publicKeySecretName. (defaults to active namespace)
- **privateKeySecretName**: Name of a new secret to add the generated private key to. The name will be generated if not specified. The secret uses format of execute-in-vm task.
- **privateKeySecretNamespace**: Namespace of privateKeySecretName. (defaults to active namespace)
- **privateKeyConnectionOptions**: Additional options to use in SSH client. Please see execute-in-vm task SSH section for more details. Eg `["host-public-key:ssh-rsa AAAAB...", "additional-ssh-options:-p 8022"]`.
- **additionalSSHKeygenOptions**: Additional options to pass to the ssh-keygen command.

### Results

- **publicKeySecretName**: The name of a public key secret.
- **publicKeySecretNamespace**: The namespace of a public key secret.
- **privateKeySecretName**: The name of a private key secret.
- **privateKeySecretNamespace**: The namespace of a private key secret.

### Usage

Task run using resolver:
```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
    generateName: generate-ssh-keys-advanced-taskrun-resolver-
spec:
    params:
    -   name: publicKeySecretName
        value: my-client-public-secret
    -   name: privateKeySecretName
        value: my-client-private-secret
    -   name: privateKeyConnectionOptions
        value:
        - user:root
        - disable-strict-host-key-checking:true
        - additional-ssh-options:-p 8022
    -   name: additionalSSHKeygenOptions
        value: -t rsa-sha2-512 -b 4096
    taskRef:
        params:
        -   name: catalog
            value: kubevirt-tekton-tasks
        -   name: type
            value: artifact
        -   name: kind
            value: task
        -   name: name
            value: generate-ssh-keys
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
    name: generate-ssh-keys-task
rules:
-   apiGroups:
    - ''
    resources:
    - secrets
    verbs:
    - get
    - list
    - create
    - patch
---
apiVersion: v1
kind: ServiceAccount
metadata:
    name: generate-ssh-keys-task
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
    name: generate-ssh-keys-task
roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: generate-ssh-keys-task
subjects:
-   kind: ServiceAccount
    name: generate-ssh-keys-task
---
```

### Platforms

The Task can be run on linux/amd64 platform.
