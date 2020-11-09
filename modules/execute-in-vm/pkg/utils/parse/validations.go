package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"strings"
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
	if strings.TrimSpace(command) == "" {
		return zerrors.NewMissingRequiredError("%v|%v option is required", commandArgsOptionName, scriptOptionName)
	}

	c.Script = command

	return nil

}
