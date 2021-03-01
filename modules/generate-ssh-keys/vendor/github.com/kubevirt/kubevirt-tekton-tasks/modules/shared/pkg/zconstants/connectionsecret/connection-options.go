package connectionsecret

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	// ConnectionSecretTypeKey supports the following values: ssh
	ConnectionSecretTypeKey = "type"
)

type sshConnectionSecretKeys struct {
	User                         string
	PrivateKey                   string
	PrivateKeyAlternativeFormat  string
	HostPublicKey                string
	DisableStrictHostKeyChecking string
	AdditionalSSHOptions         string
}

var SSHConnectionSecretKeys = sshConnectionSecretKeys{
	User:                         "user",
	PrivateKey:                   corev1.SSHAuthPrivateKey,
	PrivateKeyAlternativeFormat:  "ssh-private-key",
	HostPublicKey:                "host-public-key",
	DisableStrictHostKeyChecking: "disable-strict-host-key-checking",
	AdditionalSSHOptions:         "additional-ssh-options",
}
