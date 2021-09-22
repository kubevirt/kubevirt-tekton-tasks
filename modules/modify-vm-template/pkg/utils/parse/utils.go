package parse

import (
	"errors"
	"strings"
)

func createMapFromSlice(input []string) (map[string]string, error) {
	m := make(map[string]string)
	for _, keyValue := range input {
		key, value, err := splitPrefix(keyValue)
		if err != nil {
			return nil, err
		}
		m[key] = value
	}
	return m, nil
}

func splitPrefix(input string) (string, string, error) {
	splittedString := strings.Split(input, colonSeparator)
	if len(splittedString) < 2 {
		return "", "", errors.New("label doesn't contain : separator")
	}
	key := strings.TrimSpace(splittedString[0])
	value := strings.TrimSpace(splittedString[1])

	return key, value, nil
}
