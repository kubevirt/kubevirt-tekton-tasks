package utils

import (
	"io/ioutil"
	"os"
)

func ExitWithError(exitCode int, err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(exitCode)
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
