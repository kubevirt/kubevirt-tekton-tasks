package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"go.uber.org/zap/zapcore"
)

const (
	customizeCommandsOptionName = "customize-commands"
	customizeCommandsEnvVarName = "CUSTOMIZE_COMMANDS"
)

type CLIOptions struct {
	CustomizeCommands              string `arg:"--customize-commands,env:CUSTOMIZE_COMMANDS" placeholder:"COMMANDS" help:"virt-customize script in --commands-from-file format to execute on target pvc."`
	AdditionalVirtCustomizeOptions string `arg:"--additional-virt-customize-options,env:ADDITIONAL_VIRT_CUSTOMIZE_OPTIONS" placeholder:"OPTIONS" help:"additional options to pass to virt-customize."`
	Verbose                        string `arg:"--verbose" placeholder:"true|false" help:"Enable verbose mode and tracing of libguestfs API calls."`
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

func (c *CLIOptions) GetCustomizeCommands() string {
	return c.CustomizeCommands
}

func (c *CLIOptions) GetAdditionalVirtCustomizeOptions() string {
	return c.AdditionalVirtCustomizeOptions
}

func (c *CLIOptions) Init() error {
	if err := c.validateCommands(); err != nil {
		return err
	}

	return nil
}
