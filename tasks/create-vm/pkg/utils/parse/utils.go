package parse

import errors2 "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"

func RequireStringArgs(requiredStrings map[string]string) error {
	var requiredErrors errors2.MultiError
	for key, val := range requiredStrings {
		if val == "" {
			keyWithPrefix := "-" + key
			requiredErrors.Add(keyWithPrefix, errors2.NewMissingRequiredError("missing required argument %v", keyWithPrefix))
		}
	}

	return requiredErrors.ShortPrint("missing required arguments:").AsOptional()
}
