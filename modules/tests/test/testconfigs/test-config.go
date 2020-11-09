package testconfigs

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type TaskRunTestConfig struct {
	ServiceAccount string
	Timeout        *metav1.Duration
	LimitScope     constants.TestScope
	ExpectedLogs   string
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
