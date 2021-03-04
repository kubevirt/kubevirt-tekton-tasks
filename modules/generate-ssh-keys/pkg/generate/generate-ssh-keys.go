package generate

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/types"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func GenerateSshKeys(clioptions parse.CLIOptions) (*types.SshKeys, error) {
	opts, err := options.NewCommandOptions(clioptions.GetSshKeygenOptions())
	if err != nil {
		return nil, err
	}

	tempDir, err := ioutil.TempDir("", "sshkey-")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempDir)

	privateKeyFilename := path.Join(tempDir, "id_rsa")
	publicKeyFilename := path.Join(tempDir, "id_rsa.pub")

	setDefaultOptions(opts)
	ensureComment(opts, &clioptions)
	opts.AddOption("-f", privateKeyFilename)

	log.Logger().Debug("executing ssh-keygen command with options: " + opts.ToString())
	cmd := exec.Command(constants.SshKeyGenExecutableName, opts.GetAll()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, exit.Exit{
				Code: exitErr.ExitCode(),
				Soft: true,
			}
		} else {
			return nil, err
		}
	}

	return readKeysFromFiles(privateKeyFilename, publicKeyFilename)
}

func readKeysFromFiles(privateKeyFilename, publicKeyFilename string) (*types.SshKeys, error) {
	result := &types.SshKeys{}

	for filename, result := range map[string]*string{
		privateKeyFilename: &result.PrivateKey,
		publicKeyFilename:  &result.PublicKey,
	} {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		*result = string(content)
	}
	return result, nil
}
