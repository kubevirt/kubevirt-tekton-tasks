package env

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
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

func GetTektonResultsDir() string {
	return tektonResultsDirPath
}

func EnvOrDefault(envName string, defVal string) string {
	val, set := os.LookupEnv(envName)
	if set {
		return val
	}
	return defVal
}
