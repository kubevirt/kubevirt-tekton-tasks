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

Please see [examples](examples)

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
