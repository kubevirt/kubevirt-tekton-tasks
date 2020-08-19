package parse

import "strings"

func (c *CLIOptions) getMissingNamespaceOptionNames() string {
	var result = make([]string, 0, 2)
	if c.GetTemplateNamespace() == "" {
		result = append(result, templateNamespaceOptionName)
	}
	if c.GetVirtualMachineNamespace() == "" {
		result = append(result, vmNamespaceOptionName)
	}

	return strings.Join(result, "/")
}
