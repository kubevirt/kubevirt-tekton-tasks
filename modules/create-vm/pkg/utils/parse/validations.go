package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	"sigs.k8s.io/yaml"
	"strings"
)

func (c *CLIOptions) assertValidMode() error {
	if c.VirtualMachineManifest != "" {
		if c.TemplateName != "" {
			return zerrors.NewSoftError("only one of %v, %v should be specified", vmManifestOptionName, templateNameOptionName)
		}

		if len(c.GetTemplateParams()) > 0 || c.GetTemplateNamespace() != "" {
			return zerrors.NewSoftError("%v, %v options are not applicable for %v", templateNamespaceOptionName, templateParamsOptionName, vmManifestOptionName)
		}

	} else if c.TemplateName == "" {
		return zerrors.NewSoftError("one of %v, %v should be specified", vmManifestOptionName, templateNameOptionName)
	}

	if c.GetCreationMode() == "" {
		return zerrors.NewSoftError("could not detect correct creation mode from these options: %v, %v", vmManifestOptionName, templateNameOptionName)
	}
	return nil
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

	for _, sliceVariablePtr := range []*[]string{&c.DataVolumes, &c.OwnDataVolumes, &c.PersistentVolumeClaims, &c.OwnPersistentVolumeClaims} {
		for i, v := range *sliceVariablePtr {
			(*sliceVariablePtr)[i] = strings.TrimSpace(v)
		}
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
