package parse

import (
	"strings"
)

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

func removeVolumePrefixes(input []string) (result []string) {
	for _, keyVal := range input {
		_, val := splitVolumePrefix(keyVal)
		result = append(result, val)
	}
	return
}

func getDiskNameMap(mappings []string) map[string]string {
	result := make(map[string]string, len(mappings))

	for _, pvc := range mappings {
		key, value := splitVolumePrefix(pvc)
		if key == "" {
			key = value
		}
		result[key] = value
	}

	return result
}

func splitVolumePrefix(input string) (string, string) {
	split := strings.SplitN(input, volumesSep, 2)

	if len(split) == 2 {
		key := strings.TrimSpace(split[0])
		value := strings.TrimSpace(split[1])

		return key, value
	}

	return "", strings.TrimSpace(input)
}
