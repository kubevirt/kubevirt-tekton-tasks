package parse

import (
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	kubevirtv1 "kubevirt.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func (c *CLIOptions) assertValidMode() error {
	if c.VirtualMachineManifest != "" && c.TemplateName == "" && c.Virtctl == "" {
		if len(c.GetTemplateParams()) > 0 || c.GetTemplateNamespace() != "" {
			return zerrors.NewSoftError("%v, %v options are not applicable for %v", templateNamespaceOptionName, templateParamsOptionName, vmManifestOptionName)
		}
		return nil
	}

	if (c.VirtualMachineManifest == "" && c.TemplateName != "" && c.Virtctl == "") ||
		(c.VirtualMachineManifest == "" && c.TemplateName == "" && c.Virtctl != "") {
		return nil
	}

	return zerrors.NewSoftError("only one of %v, %v or %v should be specified", vmManifestOptionName, templateNameOptionName, virtctlOptionName)
}

func (c *CLIOptions) assertValidTypes() error {
	if !output.IsOutputType(string(c.Output)) {
		return zerrors.NewMissingRequiredError("%v is not a valid output type", c.Output)
	}
	return nil
}

func (c *CLIOptions) trimSpaces() {
	for _, strVariablePtr := range []*string{&c.TemplateName, &c.TemplateNamespace, &c.VirtualMachineNamespace} {
		*strVariablePtr = strings.TrimSpace(*strVariablePtr)
	}
}

func (c *CLIOptions) resolveDefaultNamespacesAndManifests() error {
	if c.GetCreationMode() == constants.TemplateCreationMode {
		vmNamespace := c.GetVirtualMachineNamespace()
		tempNamespace := c.GetTemplateNamespace()
		if vmNamespace == "" || tempNamespace == "" {
			activeNamespace, err := env.GetActiveNamespace()
			if err != nil {
				return zerrors.NewMissingRequiredError("%v: %v option is empty", err.Error(), c.getMissingNamespaceOptionNames())
			}
			if tempNamespace == "" {
				c.TemplateNamespace = activeNamespace
			}
			if vmNamespace == "" {
				c.VirtualMachineNamespace = activeNamespace
			}
		}
	} else if c.GetCreationMode() == constants.VMManifestCreationMode {
		vmNamespace := c.GetVirtualMachineNamespace()
		if vmNamespace == "" {
			var vm kubevirtv1.VirtualMachine

			if err := yaml.Unmarshal([]byte(c.VirtualMachineManifest), &vm); err != nil {
				return zerrors.NewMissingRequiredError("could not read VM manifest: %v", err.Error())
			}
			if vm.Namespace != "" {
				c.VirtualMachineNamespace = vm.Namespace
			} else {
				activeNamespace, err := env.GetActiveNamespace()
				if err != nil {
					return zerrors.NewMissingRequiredError("%v: %v option is empty", err.Error(), vmNamespaceOptionName)
				}
				c.VirtualMachineNamespace = activeNamespace
			}
		}
	}

	return nil
}
