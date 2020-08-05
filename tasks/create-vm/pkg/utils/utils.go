package utils

import (
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"
	"os"
)

func ErrorExit(exitCode int, err error) {
	errMsg := err.Error()
	if len(errMsg) > 0 && errMsg[len(errMsg)-1] != '\n' {
		errMsg += "\n"
	}
	_, _ = os.Stderr.WriteString(errMsg)
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
