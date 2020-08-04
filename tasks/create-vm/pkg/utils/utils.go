package utils

import (
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"
	"os"
)

func ErrorExit(exitCode int, err error) {
	_, _ = os.Stderr.WriteString(err.Error())
	os.Exit(exitCode)
}

func ErrorExitOrDie(exitCode int, err error, isSoftConditions ...bool) {
	soft := errors2.IsErrorSoft(err)

	// find any soft condition
	for idx := 0; !soft && idx < len(isSoftConditions); idx++ {
		soft = isSoftConditions[idx]
	}

	if soft {
		ErrorExit(exitCode, err)
		// end
	}
	panic(err)
}

func ConcatStringArrays(a []string, b []string) []string {
	return append(append([]string{}, a...), b...)
}
