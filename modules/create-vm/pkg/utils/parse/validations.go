package parse

import (
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	"sigs.k8s.io/yaml"
	"strings"
	"unicode"
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

func (c *CLIOptions) resolveTemplateParams() error {
	var paramsError zerrors.MultiError

	for i, param := range c.TemplateParams {
		trimmedParam := strings.TrimLeftFunc(param, unicode.IsSpace)
		c.TemplateParams[i] = trimmedParam
		split := strings.SplitN(trimmedParam, templateParamSep, 2)
		if len(split) < 2 || split[0] == "" {
			paramsError.Add(fmt.Sprintf("param %d \"%v\"", i, param), zerrors.NewMissingRequiredError("param %v has incorrect format: should be KEY:VAL", param))
		}
	}

	return paramsError.
		ShortPrint("following parameters have incorrect format: should be KEY:VAL :").
		AsOptional()
}

func (c *CLIOptions) trimSpacesAndReduceCount() {
	c.TemplateName = strings.TrimSpace(c.TemplateName)
	c.setTemplateNamespace(strings.TrimSpace(c.GetTemplateNamespace()))             // reduce count to 1
	c.setVirtualMachineNamespace(strings.TrimSpace(c.GetVirtualMachineNamespace())) // reduce count to 1

	for i, v := range c.TemplateParams {
		c.TemplateParams[i] = strings.TrimLeftFunc(v, unicode.IsSpace)
	}
	for i, v := range c.DataVolumes {
		c.DataVolumes[i] = strings.TrimSpace(v)
	}
	for i, v := range c.OwnDataVolumes {
		c.OwnDataVolumes[i] = strings.TrimSpace(v)
	}
	for i, v := range c.PersistentVolumeClaims {
		c.PersistentVolumeClaims[i] = strings.TrimSpace(v)
	}
	for i, v := range c.OwnPersistentVolumeClaims {
		c.OwnPersistentVolumeClaims[i] = strings.TrimSpace(v)
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
				c.setTemplateNamespace(activeNamespace)
			}
			if vmNamespace == "" {
				c.setVirtualMachineNamespace(activeNamespace)
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
				c.setVirtualMachineNamespace(vm.Namespace)
			} else {
				activeNamespace, err := env.GetActiveNamespace()
				if err != nil {
					return zerrors.NewMissingRequiredError("%v: %v option is empty", err.Error(), vmNamespaceOptionName)
				}
				c.setVirtualMachineNamespace(activeNamespace)
			}
		}
	}

	return nil
}
