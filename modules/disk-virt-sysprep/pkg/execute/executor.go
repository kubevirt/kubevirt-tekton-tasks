package execute

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-sysprep/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-sysprep/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-sysprep/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/options"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	cliOptions    *parse.CLIOptions
	diskImagePath string
}

func NewExecutor(clioptions *parse.CLIOptions, diskImagePath string) *Executor {
	return &Executor{cliOptions: clioptions, diskImagePath: diskImagePath}
}

func (e *Executor) PrepareGuestFSAppliance() error {
	applianceArchivePath := constants.GuestFSApplianceArchivePath

	if _, err := os.Stat(applianceArchivePath); os.IsNotExist(err) {
		return zerrors.NewMissingRequiredError("guestfs appliance is missing at %v", applianceArchivePath)
	}

	opts := []string{
		"-Jxf",
		applianceArchivePath,
		"-C",
		"/mnt",
	}

	log.GetLogger().Debug("extracting guestfs appliance with tar " + strings.Join(opts, " "))
	cmd := exec.Command("tar", opts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return os.Remove(applianceArchivePath)
}

func (e *Executor) Execute() error {
	virtSysprepScriptFileName, err := writeToTmpFile(e.cliOptions.GetSysprepCommands())
	if err != nil {
		return err
	}

	opts := options.NewCommandOptionsFromArray([]string{
		"--add",
		e.diskImagePath,
		"--commands-from-file",
		virtSysprepScriptFileName,
	})

	additionalVirtSysprepOpts, err := options.NewCommandOptions(e.cliOptions.GetAdditionalVirtSysprepOptions())
	if err != nil {
		return err
	}
	opts.AddOptions(additionalVirtSysprepOpts.GetAll()...)
	SetupVirtSysprepOptions(opts, e.cliOptions)

	log.GetLogger().Debug("executing virt-sysprep command with options: " + strings.Join(opts.GetAll(), " "))
	cmd := exec.Command("virt-sysprep", opts.GetAll()...)
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
	f, err := ioutil.TempFile("", constants.VirtSysprepCommandsFileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write([]byte(content)); err != nil {
		return "", err
	}
	return f.Name(), f.Sync()
}
