package constants

import (
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func IsEnvVarTrue(envVarName string) bool {
	return os.Getenv(envVarName) == "true"
}

func GetActiveNamespace() (string, error) {
	if activeNamespace := os.Getenv(PodNamespaceENV); activeNamespace != "" {
		return activeNamespace, nil
	}

	activeNamespaceBytes, _ := ioutil.ReadFile(serviceAccountNamespacePath)
	activeNamespace := string(activeNamespaceBytes)

	if activeNamespace != "" {
		return activeNamespace, nil
	}

	return "", errors2.NewNotFoundError("could not detect active namespace")
}

func GetTektonResultsDir() string {
	if IsEnvVarTrue(OutOfClusterENV) {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		return filepath.Dir(ex)
	}
	return TektonResultsDirPath
}
