package constants

import (
	"strings"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
)

const e2eNamespacePrefix = "e2e-tests"
const SpacesSmall = "  "

const (
	PollInterval = 1 * time.Second
)

type TargetNamespace string

const (
	DeployTargetNS TargetNamespace = "deploy"
	SystemTargetNS TargetNamespace = "system"
	EmptyTargetNS  TargetNamespace = "empty"
)

type EnvScope string

const (
	OKDEnvScope        EnvScope = "okd"
	KubernetesEnvScope EnvScope = "kubernetes"
)

func E2ETestsRandomName(name string) string {
	return strings.Join([]string{e2eNamespacePrefix, testconstants.TestRandomName(name)}, "-")
}

func E2ETestsName(name string) string {
	return strings.Join([]string{e2eNamespacePrefix, name}, "-")
}

func ToStringBoolean(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
