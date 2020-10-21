package utils

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TaskRunTestConfig struct {
	testConfig     *TestConfig
	ServiceAccount string
	Namespace      constants.TargetNamespace
	Timeout        *metav1.Duration
	LimitScope     TestScope
	ExpectedLogs   string
}

func (t *TaskRunTestConfig) GetTestConfig() *TestConfig {
	return t.testConfig
}

func (t *TaskRunTestConfig) SetTestConfig(testConfig *TestConfig) {
	t.testConfig = testConfig
}
