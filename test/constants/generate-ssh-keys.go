package constants

const (
	GenerateSshKeysTaskName = "generate-ssh-keys"
)

type generateSshKeysParams struct {
	PublicKeySecretName         string
	PublicKeySecretNamespace    string
	PrivateKeySecretName        string
	PrivateKeySecretNamespace   string
	PrivateKeyConnectionOptions string
	AdditionalSSHKeygenOptions  string
}

var GenerateSshKeysParams = generateSshKeysParams{
	PublicKeySecretName:         "publicKeySecretName",
	PublicKeySecretNamespace:    "publicKeySecretNamespace",
	PrivateKeySecretName:        "privateKeySecretName",
	PrivateKeySecretNamespace:   "privateKeySecretNamespace",
	PrivateKeyConnectionOptions: "privateKeyConnectionOptions",
	AdditionalSSHKeygenOptions:  "additionalSSHKeygenOptions",
}

type generateSshKeysResults struct {
	PublicKeySecretName       string
	PublicKeySecretNamespace  string
	PrivateKeySecretName      string
	PrivateKeySecretNamespace string
}

var GenerateSshKeysResults = generateSshKeysResults{
	PublicKeySecretName:       "publicKeySecretName",
	PublicKeySecretNamespace:  "publicKeySecretNamespace",
	PrivateKeySecretName:      "privateKeySecretName",
	PrivateKeySecretNamespace: "privateKeySecretNamespace",
}

type privateKeyConnectionOptions struct {
	Type                             string
	User                             string
	PrivateKey                       string
	HostPublicKey                    string
	DisableStrictHostKeyCheckingAttr string
	AdditionalSSHOptionsAttr         string
}

var PrivateKeyConnectionOptions = privateKeyConnectionOptions{
	Type:                             "type",
	User:                             "user",
	PrivateKey:                       "ssh-privatekey",
	HostPublicKey:                    "host-public-key",
	DisableStrictHostKeyCheckingAttr: "disable-strict-host-key-checking",
	AdditionalSSHOptionsAttr:         "additional-ssh-options",
}

const ExpectedGenerateSshKeysMessage = "The key fingerprint is"
