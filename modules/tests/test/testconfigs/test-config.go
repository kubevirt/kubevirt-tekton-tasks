package testconfigs

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type TaskRunExpectedTermination struct {
	ExitCode int32
}

type TaskRunTestConfig struct {
	ServiceAccount      string
	Timeout             *metav1.Duration
	LimitScope          constants.TestScope
	ExpectSuccess       bool
	ExpectedLogs        string
	ExpectedTermination *TaskRunExpectedTermination
}

func (t *TaskRunTestConfig) GetTaskRunTimeout() time.Duration {
	if t.Timeout != nil {
		return t.Timeout.Duration
	}
	return constants.Timeouts.DefaultTaskRun.Duration
}

func (t *TaskRunTestConfig) GetLimitScope() constants.TestScope {
	return t.LimitScope
}
