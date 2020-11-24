package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"strings"
	"time"
)

func (c *CLIOptions) trimSpacesAndReduceCount() {
	c.setVirtualMachineNamespace(strings.TrimSpace(c.GetVirtualMachineNamespace())) // reduce count to 1
}

func (c *CLIOptions) resolveDefaultNamespaces() error {
	vmNamespace := c.GetVirtualMachineNamespace()

	if vmNamespace == "" {
		activeNamespace, err := env.GetActiveNamespace()
		if err != nil {
			return zerrors.NewMissingRequiredError("%v: %v option is empty", err.Error(), vmNamespaceOptionName)
		}
		if vmNamespace == "" {
			c.setVirtualMachineNamespace(activeNamespace)
		}
	}
	return nil
}

func (c *CLIOptions) resolveExecutionScript() error {
	command := strings.Join(c.Command, " ")

	if c.GetScript() != "" {
		if command != "" {
			return zerrors.NewMissingRequiredError("only one of %v|%v options is allowed", commandOptionName, scriptOptionName)
		}
		return nil
	}
	if !c.ShouldStop() && !c.ShouldDelete() && strings.TrimSpace(command) == "" {
		return zerrors.NewMissingRequiredError("no action was specified: at least one of the following options is required: %v|%v|%v|%v",
			commandOptionName, scriptOptionName, stopOptionName, deleteOptionName)
	}

	c.Script = command

	return nil

}

func (c *CLIOptions) validateTimeout() error {
	if c.Timeout != "" {
		_, err := time.ParseDuration(c.Timeout)
		if err != nil {
			return zerrors.NewSoftError("could not parse timeout: %v", err)
		}
	}
	return nil

}

func (c *CLIOptions) validateValues() error {
	allowedValues := map[string]bool{
		"":               true,
		zconstants.False: true,
		zconstants.True:  true,
	}

	if !allowedValues[c.Stop] {
		return zerrors.NewSoftError("invalid option stop %v, only true|false is allowed", c.Stop)
	}

	if !allowedValues[c.Delete] {
		return zerrors.NewSoftError("invalid option delete %v, only true|false is allowed", c.Delete)
	}

	return nil

}
