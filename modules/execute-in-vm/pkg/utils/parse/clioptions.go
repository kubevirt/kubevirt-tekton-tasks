package parse

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"go.uber.org/zap/zapcore"
)

const (
	vmNamespaceOptionName = "vm-namespace"
	commandOptionName     = "command"
	commandArgsOptionName = "command-args"
	scriptOptionName      = "script"
)

// VirtualMachineNamespaces and Script: arrays allow to have these options without option argument
type CLIOptions struct {
	VirtualMachineName       string   `arg:"--vm-name,required" placeholder:"NAME" help:"Name of a VM to execute the action in"`
	VirtualMachineNamespaces []string `arg:"--vm-namespace" placeholder:"NAMESPACE" help:"Namespace of a VM to execute the action in"`
	Script                   string   `arg:"--script,env:EXECUTE_SCRIPT" placeholder:"SCRIPT" help:"Script to execute in a VM (can be set by EXECUTE_SCRIPT env variable)"`
	Debug                    bool     `arg:"--debug" help:"Sets DEBUG log level"`
	Command                  []string `arg:"positional" placeholder:"COMMAND" help:"Command to execute in a VM"`
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) GetVirtualMachineNamespace() string {
	return zutils.GetLast(c.VirtualMachineNamespaces)
}

func (c *CLIOptions) GetScript() string {
	return c.Script
}

func (c *CLIOptions) setVirtualMachineNamespace(namespace string) {
	if namespace == "" {
		c.VirtualMachineNamespaces = nil
	} else {
		c.VirtualMachineNamespaces = []string{namespace}
	}

}

func (c *CLIOptions) Init() error {
	c.trimSpacesAndReduceCount()

	if err := c.resolveDefaultNamespaces(); err != nil {
		return err
	}

	if err := c.resolveExecutionScript(); err != nil {
		return err
	}

	return nil
}
