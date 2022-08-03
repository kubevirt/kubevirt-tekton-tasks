package zutils

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"strings"
)

func ExtractKeysAndValuesByLastKnownKey(input []string, separator string) (map[string]string, error) {
	result := make(map[string]string, len(input))

	lastKey := ""

	for _, keyVal := range input {
		if keyVal == "" {
			continue
		}

		split := strings.SplitN(keyVal, separator, 2)

		switch len(split) {
		case 1:
			if lastKey == "" {
				return nil, zerrors.NewMissingRequiredError("no key found before \"%v\"; pair should be in \"KEY%vVAL\" format", keyVal, separator)
			} else {
				// expect space between values and append to the last key seen
				result[lastKey] += " " + split[0]

			}
		case 2:
			// key should be trimmed
			key := strings.TrimSpace(split[0])
			if key == "" {
				return nil, zerrors.NewMissingRequiredError("no key found before \"%v\"; pair should be in \"KEY%vVAL\" format", keyVal, separator)
			}
			result[key] = split[1]
			lastKey = key
		}
	}
	return result, nil
}
