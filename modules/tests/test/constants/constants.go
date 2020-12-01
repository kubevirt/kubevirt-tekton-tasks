package constants

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	"strings"
	"time"
)

const e2eNamespacePrefix = "e2e-tests"

const (
	PollInterval = 1 * time.Second
)

type TargetNamespace string

const (
	DeployTargetNS TargetNamespace = "deploy"
	TestTargetNS   TargetNamespace = "test"
	SystemTargetNS TargetNamespace = "system"
)

type TestScope string

const (
	ClusterScope   TestScope = "cluster"
	NamespaceScope TestScope = "namespace"
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
