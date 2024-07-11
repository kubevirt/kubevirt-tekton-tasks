package execute

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
)

type Executor struct {
	cliOptions    *parse.CLIOptions
	diskImagePath string
	command       string
}

func NewExecutor(clioptions *parse.CLIOptions, diskImagePath, command string) *Executor {
	return &Executor{cliOptions: clioptions, diskImagePath: diskImagePath, command: command}
}

func (e *Executor) PrepareGuestFSAppliance() error {
	applianceArchivePath := env.EnvOrDefault(constants.GuestFSApplianceArchivePathEnv, constants.GuestFSApplianceArchivePath)

	if _, err := os.Stat(applianceArchivePath); os.IsNotExist(err) {
		return zerrors.NewMissingRequiredError("guestfs appliance is missing at %v", applianceArchivePath)
	}

	os.Setenv("LIBGUESTFS_PATH", applianceArchivePath)

	return nil
}

func (e *Executor) Execute() error {
	virtScriptFileName, err := writeToTmpFile(e.cliOptions.GetCommands())
	if err != nil {
		return err
	}

	opts := options.NewCommandOptionsFromArray([]string{
		"--add",
		e.diskImagePath,
		"--commands-from-file",
		virtScriptFileName,
	})

	additionalVirtOpts, err := options.NewCommandOptions(e.cliOptions.GetAdditionalVirtOptions())
	if err != nil {
		return err
	}
	opts.AddOptions(additionalVirtOpts.GetAll()...)
	SetupVirtOptions(opts, e.cliOptions)

	log.GetLogger().Debug("executing virt command with options: " + strings.Join(opts.GetAll(), " "))
	cmd := exec.Command(e.command, opts.GetAll()...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exit.Exit{
				Code: exitErr.ExitCode(),
				Soft: true,
			}
		} else {
			return err
		}
	}
	return nil
}

func writeToTmpFile(content string) (string, error) {
	f, err := ioutil.TempFile("", constants.VirtCommandsFileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write([]byte(content)); err != nil {
		return "", err
	}
	return f.Name(), f.Sync()
}
