# Generate SSH Keys Task 

This task uses `ssh-keygen` to generate a private and public key pair

### Service Account

This task should be run with `generate-ssh-keys-task` serviceAccount.

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
