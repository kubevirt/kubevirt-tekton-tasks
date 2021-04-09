package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/requirements"
	"k8s.io/apimachinery/pkg/util/validation"
	"strings"
)

func (c *CLIOptions) trimSpaces() {
	for _, strVariablePtr := range []*string{&c.VirtualMachineInstanceName, &c.VirtualMachineInstanceNamespace, &c.SuccessCondition, &c.FailureCondition} {
		*strVariablePtr = strings.TrimSpace(*strVariablePtr)
	}
}

func (c *CLIOptions) validateNames() error {
	if c.VirtualMachineInstanceName == "" {
		return zerrors.NewMissingRequiredError("%v should not be empty", vmiNameOptionName)
	}

	for optionName, optionValue := range map[string]string{
		vmiNameOptionName:      c.VirtualMachineInstanceName,
		vmiNamespaceOptionName: c.VirtualMachineInstanceNamespace,
	} {
		if optionValue != "" {
			if errors := validation.IsDNS1123Subdomain(optionValue); len(errors) > 0 {
				return zerrors.NewMissingRequiredError("invalid %v value: %v", optionName, strings.Join(errors, ", "))
			}
		}
	}
	return nil
}

func (c *CLIOptions) resolveDefaultNamespaces() error {
	if c.VirtualMachineInstanceNamespace == "" {
		activeNamespace, err := env.GetActiveNamespace()
		if err != nil {
			return zerrors.NewMissingRequiredError("%v: %v option is empty", err.Error(), vmiNamespaceOptionName)
		}
		c.VirtualMachineInstanceNamespace = activeNamespace
	}
	return nil
}

func (c *CLIOptions) validateConditions() error {
	for conditionName, condition := range map[string]string{
		successConditionOptionName: c.SuccessCondition,
		failureConditionOptionName: c.FailureCondition,
	} {
		_, err := requirements.GetLabelRequirement(condition)
		if err != nil {
			return zerrors.NewMissingRequiredError("%v: %v", conditionName, err.Error())
		}
	}
	return nil
}
