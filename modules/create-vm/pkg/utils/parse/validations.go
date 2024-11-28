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
	if (c.VirtualMachineManifest == "" && c.Virtctl != "") || (c.VirtualMachineManifest != "" && c.Virtctl == "") {
		return nil
	}

	return zerrors.NewSoftError("only one of %v or %v should be specified", vmManifestOptionName, virtctlOptionName)
}

func (c *CLIOptions) assertValidTypes() error {
	if !output.IsOutputType(string(c.Output)) {
		return zerrors.NewMissingRequiredError("%v is not a valid output type", c.Output)
	}
	return nil
}

func (c *CLIOptions) trimSpaces() {
	c.VirtualMachineNamespace = strings.TrimSpace(c.VirtualMachineNamespace)
}

func (c *CLIOptions) resolveDefaultNamespacesAndManifests() error {
	if c.GetCreationMode() == constants.VMManifestCreationMode {
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
