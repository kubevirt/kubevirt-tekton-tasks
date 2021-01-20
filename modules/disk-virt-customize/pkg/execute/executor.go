package execute

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	clioptions    *parse.CLIOptions
	diskImagePath string
}

func NewExecutor(clioptions *parse.CLIOptions, diskImagePath string) *Executor {
	return &Executor{clioptions: clioptions, diskImagePath: diskImagePath}
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
	virtCustomizeScriptFileName, err := writeToTmpFile(e.clioptions.GetCustomizeCommands())
	if err != nil {
		return err
	}

	opts := []string{
		"--add",
		e.diskImagePath,
		"--commands-from-file",
		virtCustomizeScriptFileName,
	}

	if additionalOpts := e.clioptions.GetAdditionalVirtCustomizeOptions(); additionalOpts != "" {
		for _, additionalOpt := range strings.Fields(additionalOpts) {
			if additionalOpt != "" {
				opts = append(opts, additionalOpt)
			}
		}
	}

	log.GetLogger().Debug("executing virt-customize command with options: " + strings.Join(opts, " "))
	cmd := exec.Command("virt-customize", opts...)
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
	f, err := ioutil.TempFile("", constants.VirtCustomizeCommandsFileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write([]byte(content)); err != nil {
		return "", err
	}
	return f.Name(), f.Sync()
}
