package parse

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/logger"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"go.uber.org/zap"
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
	Command                  []string `arg:"--command" placeholder:"command" help:"Command to execute in a VM"`
	CommandArgs              []string `arg:"--command-args" placeholder:"ARG1 ARG2" help:"Arguments of a command"`
	Scripts                  []string `arg:"--script" placeholder:"SCRIPT" help:"Script to execute in a VM"`
	Debug                    bool     `arg:"--debug" help:"Sets DEBUG log level"`
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
	return zutils.GetLast(c.Scripts)
}

func (c *CLIOptions) setVirtualMachineNamespace(namespace string) {
	if namespace == "" {
		c.VirtualMachineNamespaces = nil
	} else {
		c.VirtualMachineNamespaces = []string{namespace}
	}

}

func (c *CLIOptions) setScript(script string) {
	if script == "" {
		c.Scripts = nil
	} else {
		c.Scripts = []string{script}
	}
}

func (c *CLIOptions) Init() error {
	defer logger.GetLogger().Debug("parsed arguments", zap.Reflect("cliOptions", c))

	c.trimSpacesAndReduceCount()

	if err := c.resolveDefaultNamespaces(); err != nil {
		return err
	}

	if err := c.resolveExecutionScript(); err != nil {
		return err
	}

	return nil
}
