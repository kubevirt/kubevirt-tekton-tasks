package execattributes

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/env/fileoptions"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"strings"
)

// General options
const (
	// supported values: ExecSecretType.SSH (ssh)
	secretTypeOption = "type"
)

type attributes struct {
	secretType constants.ExecSecretType
	secretPath string
	ssh        SSHAttributes
}

func NewExecAttributes() ExecAttributes {
	return &attributes{}
}

type ExecAttributes interface {
	zapcore.ObjectMarshaler
	Init(execAttributesPath string) error
	GetType() constants.ExecSecretType
	GetSSHAttributes() SSHAttributes
}

func (s *attributes) Init(execAttributesPath string) error {
	if _, err := os.Stat(execAttributesPath); os.IsNotExist(err) {
		return zerrors.NewMissingRequiredError("secret does not exist at %v", execAttributesPath)
	}
	s.secretPath = execAttributesPath

	var secretTypeRaw string
	if err := fileoptions.ReadFileOption(&secretTypeRaw, path.Join(s.secretPath, secretTypeOption)); err != nil {
		return err
	}

	secretTypeRaw = strings.TrimSpace(secretTypeRaw)

	switch secretTypeRaw {
	case string(constants.SSHSecretType):
		s.secretType = constants.ExecSecretType(secretTypeRaw)
	default:
		if secretTypeRaw == "" {
			return zerrors.NewMissingRequiredError("%v secret attribute is required", secretTypeOption)
		}
		return zerrors.NewMissingRequiredError("%v is invalid %v", secretTypeRaw, secretTypeOption)
	}

	switch s.secretType {
	case constants.SSHSecretType:
		s.ssh = NewSSHAttributes()
		if err := s.ssh.initSSH(s.secretPath); err != nil {
			return err
		}
	}

	return nil
}

func (s *attributes) GetType() constants.ExecSecretType {
	return s.secretType
}

func (s *attributes) GetSSHAttributes() SSHAttributes {
	return s.ssh
}

func (s *attributes) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("secretType", string(s.secretType))
	encoder.AddString("secretPath", s.secretPath)
	return encoder.AddObject("ssh", s.ssh)
}
