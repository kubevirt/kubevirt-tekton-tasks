package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"k8s.io/apimachinery/pkg/util/validation"
	"strings"
	"time"
)

func (c *CLIOptions) trimSpaces() {
	c.VirtualMachineNamespace = strings.TrimSpace(c.VirtualMachineNamespace)
}

func (c *CLIOptions) validateName() error {
	if c.VirtualMachineName == "" {
		return zerrors.NewMissingRequiredError("missing value for %v option", vmNameOptionName)
	}

	errs := validation.IsDNS1123Subdomain(c.VirtualMachineName)
	if len(errs) > 0 {
		return zerrors.NewMissingRequiredError("%v is not a valid name: %v", vmNameOptionName, strings.Join(errs, ";"))
	}
	return nil
}

func (c *CLIOptions) resolveDefaultNamespaces() error {
	vmNamespace := c.GetVirtualMachineNamespace()

	if vmNamespace == "" {
		activeNamespace, err := env.GetActiveNamespace()
		if err != nil {
			return zerrors.NewMissingRequiredError("%v: %v option is empty", err.Error(), vmNamespaceOptionName)
		}
		if vmNamespace == "" {
			c.VirtualMachineNamespace = activeNamespace
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

func (c *CLIOptions) validateConnectionSecretName() error {
	command := strings.Join(c.Command, " ")

	if c.GetScript() != "" || strings.TrimSpace(command) != "" {
		if c.ConnectionSecretName == "" || c.ConnectionSecretName == constants.EmptyConnectionSecretName {
			return zerrors.NewMissingRequiredError("connection secret should not be empty")
		}

		errs := validation.IsDNS1123Subdomain(c.ConnectionSecretName)
		if len(errs) > 0 {
			return zerrors.NewMissingRequiredError("connection secret does not have a valid name: %v", strings.Join(errs, ";"))
		}
	}

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
