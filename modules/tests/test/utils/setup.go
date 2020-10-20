package utils

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type TestConfig struct {
	TestOptions
	RestConfig *rest.Config
}

func (t *TestConfig) LimitScope(limitScope TestScope) {
	if limitScope != "" && limitScope != t.Scope {
		ginkgo.Skip(fmt.Sprintf("runs only in %v scope", limitScope))
	}
}

func (t *TestConfig) GetResolvedTestNamespace(namespace constants.TargetNamespace) string {
	if namespace == constants.DeployTargetNS {
		return t.DeployNamespace
	} else if namespace == constants.CustomTargetNS {
		if t.DeployNamespace == "tekton-pipelines" {
			return "default"
		} else {
			return "tekton-pipelines"
		}
	}
	return t.TestNamespace
}

func Setup() (*TestConfig, error) {
	testOpts, err := GetTestOptions()
	if err != nil {
		return nil, err
	}

	restConf, err := rest.InClusterConfig()
	if err != nil {
		restConf, err = clientcmd.BuildConfigFromFlags("", testOpts.KubeConfigPath)
	}
	if err != nil {
		return nil, fmt.Errorf("could not load KUBECONFIG: %v", err)
	}

	return &TestConfig{
		RestConfig:  restConf,
		TestOptions: *testOpts,
	}, nil
}
