package execute

import (
	"fmt"
	cmd2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/cmd"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execattributes"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	knownHostsFilename = "known_hosts"
	idRSAFilename      = "id_rsa"
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

	log.Logger().Debug("preparing ssh files")
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
		log.Logger().Debug("connection not found: " + address)
	}

	return conn != nil && err == nil
}

func (e *sshExecutor) RemoteExecute(timeout time.Duration) error {
	opts := options.NewCommandOptionsFromArray(e.ssh.GetAdditionalSSHOptions())

	destination := e.ssh.GetUser() + "@" + e.ipAddress
	opts.AddValue(destination)

	log.Logger().Debug("executing ssh command with options: " + strings.Join(opts.GetAll(), " "))

	// do not log script
	opts.AddValue("--")
	opts.AddValue(e.clioptions.GetScript())

	cmd := exec.Command(e.ssh.GetSSHExecutableName(), opts.GetAll()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd2.RunCmdWithTimeout(timeout, cmd)
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
