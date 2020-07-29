package parse

import createVMerrs "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"

func RequireStringArgs(requiredStrings map[string]string) error {
	var requiredErrors []error
	for key, val := range requiredStrings {
		if val == "" {
			requiredErrors = append(requiredErrors, createVMerrs.NewMissingRequiredArgError(key))
		}
	}

	return createVMerrs.NewMultiErrorOrNil(requiredErrors)
}
