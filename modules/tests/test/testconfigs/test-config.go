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
	LimitTestScope      constants.TestScope
	LimitEnvScope       constants.EnvScope
	ExpectSuccess       bool
	ExpectedLogs        string
	ExpectedLogsList    []string
	ExpectedTermination *TaskRunExpectedTermination
}

func (t *TaskRunTestConfig) GetTaskRunTimeout() time.Duration {
	if t.Timeout != nil {
		return t.Timeout.Duration
	}
	return constants.Timeouts.DefaultTaskRun.Duration
}

func (t *TaskRunTestConfig) GetLimitTestScope() constants.TestScope {
	return t.LimitTestScope
}

func (t *TaskRunTestConfig) GetLimitEnvScope() constants.EnvScope {
	return t.LimitEnvScope
}

func (t *TaskRunTestConfig) GetAllExpectedLogs() []string {
	var allLogs []string

	if t.ExpectedLogs != "" {
		allLogs = append(allLogs, t.ExpectedLogs)
	}

	if t.ExpectedLogsList != nil {
		allLogs = append(allLogs, t.ExpectedLogsList...)
	}

	return allLogs
}
