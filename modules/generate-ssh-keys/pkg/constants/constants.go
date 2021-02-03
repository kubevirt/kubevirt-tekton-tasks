package constants

// Exit codes
// reserve 0+ numbers for the exit code of the command
const (
	InvalidArguments               = -1 // same as go-arg invalid args exit
	SshKeysGenerationFailed        = -2
	PrivateKeyAlreadyExists        = -3
	SecretCreatorInitFailed        = -4
	PublicKeySecretFetchFailed     = -5
	PublicKeySecretCreationFailed  = -6
	PrivateKeySecretCreationFailed = -7
	WriteResultsExitCode           = -8
)

type results struct {
	PublicKeySecretName       string
	PublicKeySecretNamespace  string
	PrivateKeySecretName      string
	PrivateKeySecretNamespace string
}

var Results = results{
	PublicKeySecretName:       "publicKeySecretName",
	PublicKeySecretNamespace:  "publicKeySecretNamespace",
	PrivateKeySecretName:      "privateKeySecretName",
	PrivateKeySecretNamespace: "privateKeySecretNamespace",
}

type connectionOptions struct {
	Type                             string
	User                             string
	PrivateKey                       string
	HostPublicKey                    string
	DisableStrictHostKeyCheckingAttr string
	AdditionalSSHOptionsAttr         string
}

var ConnectionOptions = connectionOptions{
	Type:                             "type",
	User:                             "user",
	PrivateKey:                       "private-key",
	HostPublicKey:                    "host-public-key",
	DisableStrictHostKeyCheckingAttr: "disable-strict-host-key-checking",
	AdditionalSSHOptionsAttr:         "additional-ssh-options",
}

const SshKeyGenExecutableName = "ssh-keygen"
const ConnectionSSHType = "ssh"

const (
	PrivateKeyGenerateName = "private-key-"
	PublicKeyGenerateName  = "public-key-"
)
