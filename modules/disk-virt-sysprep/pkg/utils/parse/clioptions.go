package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"go.uber.org/zap/zapcore"
)

const (
	sysprepCommandsOptionName = "sysprep-commands"
	sysprepCommandsEnvVarName = "SYSPREP_COMMANDS"
)

type CLIOptions struct {
	SysprepCommands              string `arg:"--sysprep-commands,env:SYSPREP_COMMANDS" placeholder:"COMMANDS" help:"virt-sysprep script in --commands-from-file format to execute on target pvc."`
	AdditionalVirtSysprepOptions string `arg:"--additional-virt-sysprep-options,env:ADDITIONAL_VIRT_SYSPREP_OPTIONS" placeholder:"OPTIONS" help:"additional options to pass to virt-sysprepr."`
	Verbose                      string `arg:"--verbose" placeholder:"true|false" help:"Enable verbose mode and tracing of libguestfs API calls."`
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

func (c *CLIOptions) GetSysprepCommands() string {
	return c.SysprepCommands
}

func (c *CLIOptions) GetAdditionalVirtSysprepOptions() string {
	return c.AdditionalVirtSysprepOptions
}

func (c *CLIOptions) Init() error {
	if err := c.validateCommands(); err != nil {
		return err
	}

	return nil
}
