package execute

import (
	"fmt"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execattributes"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/parse"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const (
	knownHostsFilename = "known_hosts"
	idRSAFilename      = "id_rsa"
	sshPort            = "22"
)

const (
	defaultFileMode = 0600
	defaultDirMode  = 0700
)

type sshExecutor struct {
	clioptions *parse.CLIOptions
	ssh        execattributes.SSHAttributes
	ipAddress  string
}

func newSSHExecutor(clioptions *parse.CLIOptions, execAttributes execattributes.ExecAttributes) *sshExecutor {
	return &sshExecutor{clioptions: clioptions, ssh: execAttributes.GetSSHAttributes()}
}

func (e *sshExecutor) Init(ipAddress string) error {
	e.ipAddress = ipAddress

	log.GetLogger().Debug("preparing ssh files")
	if err := os.MkdirAll(e.ssh.GetSSHDir(), defaultDirMode); err != nil {
		return err
	}

	if privateKey := e.ssh.GetPrivateKey(); privateKey != "" {
		if err := writeToUserFile(path.Join(e.ssh.GetSSHDir(), idRSAFilename), privateKey, false); err != nil {
			return err
		}
	}

	if hostPublicKey := e.ssh.GetHostPublicKey(); hostPublicKey != "" {
		knownHost := fmt.Sprintf("%v %v", ipAddress, hostPublicKey)
		if err := writeToUserFile(path.Join(e.ssh.GetSSHDir(), knownHostsFilename), knownHost, true); err != nil {
			return err
		}
	}

	return nil

}

func (e *sshExecutor) TestConnection() bool {
	address := net.JoinHostPort(e.ipAddress, strconv.Itoa(e.ssh.GetPort()))
	conn, err := net.DialTimeout("tcp", address, constants.CheckSSHConnectionTimeout)
	if conn != nil {
		defer conn.Close()
	} else {
		log.GetLogger().Debug("connection not found: " + address)
	}

	return conn != nil && err == nil
}

func (e *sshExecutor) RemoteExecute() error {
	destination := e.ssh.GetUser() + "@" + e.ipAddress

	var opts []string
	if additionalOpts := e.ssh.GetAdditionalSSHOptions(); additionalOpts != "" {
		for _, additionalOpt := range strings.Fields(additionalOpts) {
			if additionalOpt != "" {
				opts = append(opts, additionalOpt)
			}
		}
	}
	opts = append(opts, destination)
	log.GetLogger().Debug("executing ssh command with options: " + strings.Join(opts, " "))
	// do not log script
	opts = append(opts, "--")
	opts = append(opts, e.clioptions.GetScript())

	cmd := exec.Command(e.ssh.GetSSHExecutableName(), opts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	result := exit.Exit{
		Code: 0,
		Soft: true,
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.Code = exitErr.ExitCode()
		} else {
			return err
		}
	}

	return result
}

func writeToUserFile(filename string, content string, append bool) error {
	flags := os.O_CREATE | os.O_WRONLY

	if append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	f, err := os.OpenFile(filename, flags, defaultFileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte(content)); err != nil {
		return err
	}
	return f.Sync()
}
