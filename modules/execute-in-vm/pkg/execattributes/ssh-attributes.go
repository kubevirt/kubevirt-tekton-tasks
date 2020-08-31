package execattributes

import (
	"fmt"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/env/fileoptions"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap/zapcore"
	"os/user"
	"path"
	"strings"
)

const (
	sshUserAttr                         = "user"
	sshPrivateKeyAttr                   = "private-key"
	sshHostPublicKeyAttr                = "host-public-key"
	sshDisableStrictHostKeyCheckingAttr = "disable-strict-host-key-checking"
	sshAdditionalSSHOptionsAttr         = "additional-ssh-options"
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
	yes       = "yes"
)

// port is parsed from additionalSSHOptions to be more easier to use
type sshAttributes struct {
	user                         string
	port                         int
	additionalSSHOptions         string
	privateKey                   string
	hostPublicKey                string
	disableStrictHostKeyChecking bool
}

type SSHAttributes interface {
	zapcore.ObjectMarshaler
	initSSH(execSecretPath string) error
	GetUser() string
	GetPort() int
	GetAdditionalSSHOptions() string
	IncludesSSHOption(option string) bool
	AddAdditionalSSHOption(name, value string)
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

	stringOptions := map[string]*string{
		sshUserAttr:                 &s.user,
		sshAdditionalSSHOptionsAttr: &s.additionalSSHOptions,
		sshPrivateKeyAttr:           &s.privateKey,
		sshHostPublicKeyAttr:        &s.hostPublicKey,
	}
	boolOptions := map[string]*bool{
		sshDisableStrictHostKeyCheckingAttr: &s.disableStrictHostKeyChecking,
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

	if s.user == "" {
		return zerrors.NewMissingRequiredError("%v secret attribute is required", sshUserAttr)
	}

	if strings.TrimSpace(s.privateKey) == "" {
		return zerrors.NewMissingRequiredError("%v secret attribute is required", sshPrivateKeyAttr)
	}

	if strings.TrimSpace(s.hostPublicKey) == "" && !s.disableStrictHostKeyChecking {
		return zerrors.NewMissingRequiredError("%v or %v=true secret attribute is required", sshHostPublicKeyAttr, sshDisableStrictHostKeyCheckingAttr)
	}
	port, err := parsePort(s.additionalSSHOptions)
	if err != nil {
		return err
	}
	s.port = port

	if !s.IncludesSSHOption(sshStrictHostKeyCheckingOption) {
		s.AddAdditionalSSHOption(sshStrictHostKeyCheckingOption, s.GetStrictHostKeyCheckingMode())
	}

	return nil
}

func (s *sshAttributes) GetUser() string {
	return s.user
}

func (s *sshAttributes) GetPort() int {
	return s.port
}

func (s *sshAttributes) GetAdditionalSSHOptions() string {
	return s.additionalSSHOptions
}

func (s *sshAttributes) IncludesSSHOption(option string) bool {
	return strings.Contains(s.additionalSSHOptions, option)
}

func (s *sshAttributes) AddAdditionalSSHOption(name, value string) {
	if s.additionalSSHOptions != "" {
		s.additionalSSHOptions += " "
	}
	s.additionalSSHOptions += fmt.Sprintf("-o%v=%v", name, value)
}

func (s *sshAttributes) GetStrictHostKeyCheckingMode() string {
	if s.disableStrictHostKeyChecking {
		return acceptNew
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
	encoder.AddString("additionalSSHOptions", s.additionalSSHOptions)
	encoder.AddBool("disableStrictHostKeyChecking", s.disableStrictHostKeyChecking)
	return nil
}
