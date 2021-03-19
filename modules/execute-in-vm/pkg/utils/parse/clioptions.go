package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	vmNameOptionName      = "vm-name"
	vmNamespaceOptionName = "vm-namespace"
	stopOptionName        = "stop"
	deleteOptionName      = "delete"
	commandOptionName     = "command"
	commandArgsOptionName = "command-args"
	scriptOptionName      = "script"
)

type CLIOptions struct {
	VirtualMachineName      string   `arg:"--vm-name,env:VM_NAME,required" placeholder:"NAME" help:"Name of a VM to execute the action in"`
	VirtualMachineNamespace string   `arg:"--vm-namespace,env:VM_NAMESPACE" placeholder:"NAMESPACE" help:"Namespace of a VM to execute the action in"`
	Stop                    string   `arg:"--stop" placeholder:"true|false" help:"Stops the VM after executing the action"`
	Delete                  string   `arg:"--delete" placeholder:"true|false" help:"Deletes the VM after executing the action"`
	Timeout                 string   `arg:"--timeout" help:"Timeout for the command/script (includes potential VM start). The VM will be stoped or deleted accordingly once the timout expires. Should be in a 3h2m1s format."`
	Script                  string   `arg:"--script,env:EXECUTE_SCRIPT" placeholder:"SCRIPT" help:"Script to execute in a VM (can be set by EXECUTE_SCRIPT env variable)"`
	ConnectionSecretName    string   `arg:"--connectionSecretName,env:CONNECTION_SECRET_NAME" placeholder:"NAME" help:"Name of the connection secret (used only for validation)"`
	Debug                   bool     `arg:"--debug" help:"Sets DEBUG log level"`
	Command                 []string `arg:"positional" placeholder:"COMMAND" help:"Command to execute in a VM"`
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) GetVirtualMachineNamespace() string {
	return c.VirtualMachineNamespace
}

func (c *CLIOptions) GetScript() string {
	return c.Script
}

func (c *CLIOptions) GetScriptTimeout() time.Duration {
	if c.Timeout != "" {
		timeout, err := time.ParseDuration(c.Timeout)
		if err == nil {
			return timeout
		}
	}

	return 0
}

func (c *CLIOptions) ShouldStop() bool {
	return zutils.IsTrue(c.Stop)
}

func (c *CLIOptions) ShouldDelete() bool {
	return zutils.IsTrue(c.Delete)
}

func (c *CLIOptions) Init() error {
	c.trimSpaces()

	if err := c.validateName(); err != nil {
		return err
	}

	if err := c.resolveDefaultNamespaces(); err != nil {
		return err
	}

	if err := c.resolveExecutionScript(); err != nil {
		return err
	}

	if err := c.validateConnectionSecretName(); err != nil {
		return err
	}

	if err := c.validateTimeout(); err != nil {
		return err
	}

	if err := c.validateValues(); err != nil {
		return err
	}

	return nil
}
