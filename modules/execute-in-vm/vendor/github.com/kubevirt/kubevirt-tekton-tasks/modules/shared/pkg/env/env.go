package env

import (
	"errors"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	serviceAccountNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	tektonResultsDirPath        = "/tekton/results"
)

func IsEnvVarTrue(envVarName string) bool {
	return zutils.IsTrue(os.Getenv(envVarName))
}

func GetActiveNamespace() (string, error) {
	activeNamespaceBytes, err := ioutil.ReadFile(serviceAccountNamespacePath)

	if err == nil {
		if activeNamespace := string(activeNamespaceBytes); activeNamespace != "" {
			return activeNamespace, nil
		}
	}

	return "", errors.New("could not detect active namespace")
}

func GetWorkingDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func GetTektonResultsDir() string {
	if IsEnvVarTrue(zconstants.OutOfClusterENV) {
		return GetWorkingDir()
	}
	return tektonResultsDirPath
}
