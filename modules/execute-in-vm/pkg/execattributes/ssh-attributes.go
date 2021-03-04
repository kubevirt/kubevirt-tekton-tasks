package execattributes

import (
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env/fileoptions"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants/connectionsecret"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap/zapcore"
	"os/user"
	"path"
	"strings"
)

const (
	sshStrictHostKeyCheckingOption = "StrictHostKeyChecking"
)

const (
	sshDirName        = ".ssh"
	sshExecutableName = "ssh"
)

const (
	acceptNew = "accept-new"
	no        = "no"
	yes       = "yes"
)

// port is parsed from additionalSSHOptions to be more easier to use
type sshAttributes struct {
	user                         string
	port                         int
	additionalSSHOptions         []string
	privateKey                   string
	hostPublicKey                string
	disableStrictHostKeyChecking bool
}

type SSHAttributes interface {
	zapcore.ObjectMarshaler
	initSSH(execSecretPath string) error
	GetUser() string
	GetPort() int
	GetAdditionalSSHOptions() []string
	GetPrivateKey() string
	GetHostPublicKey() string
	GetStrictHostKeyCheckingMode() string
	GetSSHDir() string
	GetSSHExecutableName() string
}

func NewSSHAttributes() SSHAttributes {
	return &sshAttributes{}
}

func (s *sshAttributes) initSSH(execSecretPath string) error {
	var privateKeyAlternativeFormat, additionalSSHOptionsString string

	stringOptions := map[string]*string{
		connectionsecret.SSHConnectionSecretKeys.User:                        &s.user,
		connectionsecret.SSHConnectionSecretKeys.AdditionalSSHOptions:        &additionalSSHOptionsString,
		connectionsecret.SSHConnectionSecretKeys.PrivateKey:                  &s.privateKey,
		connectionsecret.SSHConnectionSecretKeys.PrivateKeyAlternativeFormat: &privateKeyAlternativeFormat,
		connectionsecret.SSHConnectionSecretKeys.HostPublicKey:               &s.hostPublicKey,
	}
	boolOptions := map[string]*bool{
		connectionsecret.SSHConnectionSecretKeys.DisableStrictHostKeyChecking: &s.disableStrictHostKeyChecking,
	}

	for optionName, output := range stringOptions {
		if err := fileoptions.ReadFileOption(output, path.Join(execSecretPath, optionName)); err != nil {
			return err
		}
	}

	for optionName, output := range boolOptions {
		if err := fileoptions.ReadFileOptionBool(output, path.Join(execSecretPath, optionName)); err != nil {
			return err
		}
	}

	if strings.TrimSpace(s.privateKey) == "" {
		if strings.TrimSpace(privateKeyAlternativeFormat) != "" {
			s.privateKey = privateKeyAlternativeFormat
		} else {
			return zerrors.NewMissingRequiredError("%v secret attribute is required", connectionsecret.SSHConnectionSecretKeys.PrivateKey)
		}
	}

	if s.user == "" {
		return zerrors.NewMissingRequiredError("%v secret attribute is required", connectionsecret.SSHConnectionSecretKeys.User)
	}

	if strings.TrimSpace(s.hostPublicKey) == "" && !s.disableStrictHostKeyChecking {
		return zerrors.NewMissingRequiredError("%v or %v=true secret attribute is required", connectionsecret.SSHConnectionSecretKeys.HostPublicKey, connectionsecret.SSHConnectionSecretKeys.DisableStrictHostKeyChecking)
	}
	additionalSSHOptions, err := options.NewCommandOptions(additionalSSHOptionsString)
	if err != nil {
		return err
	}

	port, err := parsePort(additionalSSHOptions)
	if err != nil {
		return err
	}
	s.port = port

	if !strings.HasSuffix(s.privateKey, "\n") {
		s.privateKey += "\n"
	}

	if !additionalSSHOptions.IncludesString(sshStrictHostKeyCheckingOption) {
		additionalSSHOptions.AddOption("-o", fmt.Sprintf("%v=%v", sshStrictHostKeyCheckingOption, s.GetStrictHostKeyCheckingMode()))
	}

	s.additionalSSHOptions = additionalSSHOptions.GetAll()

	return nil
}

func (s *sshAttributes) GetUser() string {
	return s.user
}

func (s *sshAttributes) GetPort() int {
	return s.port
}

func (s *sshAttributes) GetAdditionalSSHOptions() []string {
	return s.additionalSSHOptions
}

func (s *sshAttributes) GetStrictHostKeyCheckingMode() string {
	if s.disableStrictHostKeyChecking {
		// TODO change to safer acceptNew once a newer version of ssh which supports this option is available in CI
		return no
	}
	return yes
}

func (s *sshAttributes) GetPrivateKey() string {
	return s.privateKey
}

func (s *sshAttributes) GetHostPublicKey() string {
	return s.hostPublicKey
}

func (s *sshAttributes) GetSSHDir() string {
	current, err := user.Current()

	if err != nil {
		panic(err)
	}

	return path.Join(current.HomeDir, sshDirName)
}

func (s *sshAttributes) GetSSHExecutableName() string {
	return sshExecutableName
}

func (s *sshAttributes) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	// do not print private/public key
	encoder.AddString("user", s.user)
	encoder.AddString("additionalSSHOptions", strings.Join(s.additionalSSHOptions, " "))
	encoder.AddBool("disableStrictHostKeyChecking", s.disableStrictHostKeyChecking)
	return nil
}
