package constants

// Exit codes
// reserve 0+ numbers for the exit code of the command
const (
	InvalidArguments               = -1 // same as go-arg invalid args exit
	SshKeysGenerationFailed        = -2
	PrivateKeyAlreadyExists        = -3
	SecretFacadeInitFailed         = -4
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

const SshKeyGenExecutableName = "ssh-keygen"

const (
	PrivateKeyGenerateName = "private-key-"
	PublicKeyGenerateName  = "public-key-"
)
