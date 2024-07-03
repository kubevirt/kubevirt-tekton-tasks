package parse

import (
	"go.uber.org/zap/zapcore"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
)

const (
	commandsOptionName = "virt-commands"
	commandsEnvVarName = "VIRT_COMMANDS"
)

type CLIOptions struct {
	Commands              string `arg:"--virt-commands,env:VIRT_COMMANDS" placeholder:"VIRT_COMMANDS" help:"virt script in --commands-from-file format to execute on target pvc."`
	AdditionalVirtOptions string `arg:"--additional-virt-options,env:ADDITIONAL_VIRT_OPTIONS" placeholder:"OPTIONS" help:"additional options to pass to virt command"`
	Verbose               string `arg:"--verbose" placeholder:"true|false" help:"Enable verbose mode and tracing of libguestfs API calls."`
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.IsVerbose() {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) IsVerbose() bool {
	return zutils.IsTrue(c.Verbose)
}

func (c *CLIOptions) GetCommands() string {
	return c.Commands
}

func (c *CLIOptions) GetAdditionalVirtOptions() string {
	return c.AdditionalVirtOptions
}

func (c *CLIOptions) Init() error {
	if err := c.validateCommands(); err != nil {
		return err
	}

	return nil
}
