package parse

import (
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	publicKeySecretNameOptionName       = "public-key-secret-name"
	publicKeySecretNamespaceOptionName  = "public-key-secret-namespace"
	privateKeySecretNameOptionName      = "private-key-secret-name"
	privateKeySecretNamespaceOptionName = "private-key-secret-namespace"
)

const connectionOptionsSep = ":"

type CLIOptions struct {
	PublicKeySecretName         string   `arg:"--public-key-secret-name,env:PUBLIC_KEY_SECRET_NAME" placeholder:"NAME" help:"Name of a new or existing secret to append the generated public key to. The name will be generated and new secret created if not specified."`
	PublicKeySecretNamespace    string   `arg:"--public-key-secret-namespace,env:PUBLIC_KEY_SECRET_NAMESPACE" placeholder:"NAME" help:"Namespace of public-key-secret-name. (defaults to active namespace)"`
	PrivateKeySecretName        string   `arg:"--private-key-secret-name,env:PRIVATE_KEY_SECRET_NAME" placeholder:"NAME" help:"Name of a new secret to add the generated private key to. The name will be generated if not specified. The secret uses format of execute-in-vm task."`
	PrivateKeySecretNamespace   string   `arg:"--private-key-secret-namespace,env:PRIVATE_KEY_SECRET_NAMESPACE" placeholder:"NAME" help:"Namespace of private-key-secret-name. (defaults to active namespace)"`
	SshKeygenOptions            string   `arg:"--additional-ssh-keygen-options,env:ADDITIONAL_SSH_KEYGEN_OPTIONS" placeholder:"OPTIONS" help:"Additional options to pass to the ssh-keygen command."`
	Debug                       bool     `arg:"--debug" help:"Sets DEBUG log level"`
	PrivateKeyConnectionOptions []string `arg:"positional" placeholder:"KEY1:VAL1 KEY2:VAL2" help:"Additional private-key connection options to use in SSH client. Please see execute-in-vm task SSH section for more details. Eg [\"host-public-key:ssh-rsa AAAAB...\", \"additional-ssh-options:-p 8022\"]."`
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) GetPublicKeySecretName() string {
	return c.PublicKeySecretName
}
func (c *CLIOptions) GetPublicKeySecretNamespace() string {
	return c.PublicKeySecretNamespace
}

func (c *CLIOptions) GetPrivateKeySecretName() string {
	return c.PrivateKeySecretName
}
func (c *CLIOptions) GetPrivateKeySecretNamespace() string {
	return c.PrivateKeySecretNamespace
}

func (c *CLIOptions) GetSshKeygenOptions() string {
	return c.SshKeygenOptions
}

func (c *CLIOptions) GetPrivateKeyConnectionOptions() map[string]string {
	result := make(map[string]string, len(c.PrivateKeyConnectionOptions))

	lastKey := ""

	for _, keyVal := range c.PrivateKeyConnectionOptions {
		split := strings.SplitN(keyVal, connectionOptionsSep, 2)

		switch len(split) {
		case 1:
			// expect space between values and append to the last key seen
			if lastKey != "" {
				result[lastKey] += " " + split[0]
			}
		case 2:
			lastKey = strings.TrimSpace(split[0])
			result[lastKey] = split[1]
		}
	}
	return result
}

func (c *CLIOptions) Init() error {
	c.trimSpaces()

	if err := c.validateNames(); err != nil {
		return err
	}

	if err := c.resolveDefaultNamespaces(); err != nil {
		return err
	}
	return nil
}
