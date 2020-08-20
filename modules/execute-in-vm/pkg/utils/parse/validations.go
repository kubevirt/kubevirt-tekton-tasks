package parse

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"strings"
)

func (c *CLIOptions) trimSpacesAndReduceCount() {
	c.setScript(c.GetScript())                                                      // reduce count to 1
	c.setVirtualMachineNamespace(strings.TrimSpace(c.GetVirtualMachineNamespace())) // reduce count to 1
}

func (c *CLIOptions) resolveDefaultNamespaces() error {
	vmNamespace := c.GetVirtualMachineNamespace()

	if vmNamespace == "" {
		activeNamespace, err := zconstants.GetActiveNamespace()
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
	args := strings.Join(c.CommandArgs, " ")

	if c.GetScript() != "" {
		if command != "" {
			return zerrors.NewMissingRequiredError("only one of %v|%v options is allowed", commandOptionName, scriptOptionName)
		}
		if args != "" {
			return zerrors.NewMissingRequiredError("only one of %v|%v options is allowed", commandArgsOptionName, scriptOptionName)
		}
		return nil
	}
	if strings.TrimSpace(command) == "" {
		return zerrors.NewMissingRequiredError("one of %v|%v options is required", commandArgsOptionName, scriptOptionName)
	}

	if args != "" {
		command += " " + args
	}

	c.setScript(command)

	return nil

}
