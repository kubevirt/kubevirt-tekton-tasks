package utils

import (
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"
	"io/ioutil"
	"os"
)

func Exit(exitCode int, err error) {
	_, _ = os.Stderr.WriteString(err.Error())
	os.Exit(exitCode)
}

func ExitOrDie(exitCode int, err error, isSoftConditions ...bool) {
	soft := errors2.IsErrorSoft(err)

	// find any soft condition
	for idx := 0; !soft && idx < len(isSoftConditions); idx++ {
		soft = isSoftConditions[idx]
	}

	if soft {
		Exit(exitCode, err)
		// end
	}
	panic(err)
}

func WriteToFile(path string, content string) {
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

func ConcatStringArrays(a []string, b []string) []string {
	return append(append([]string{}, a...), b...)
}
