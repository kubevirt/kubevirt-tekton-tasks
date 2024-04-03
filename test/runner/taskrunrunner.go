package runner

import (
	"context"

	"github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	framework2 "github.com/kubevirt/kubevirt-tekton-tasks/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/tekton"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/testconfigs"
	. "github.com/onsi/gomega"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tkntest "github.com/tektoncd/pipeline/test"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

type TaskRunRunner struct {
	framework *framework2.Framework
	taskRun   *pipev1.TaskRun
	logs      string
}

func NewTaskRunRunner(framework *framework2.Framework, taskRun *pipev1.TaskRun) *TaskRunRunner {
	Expect(taskRun).ShouldNot(BeNil())
	return &TaskRunRunner{
		framework: framework,
		taskRun:   taskRun,
	}
}

func (r *TaskRunRunner) GetTaskRun() *pipev1.TaskRun {
	return r.taskRun
}

func (r *TaskRunRunner) CreateTaskRun() *TaskRunRunner {
	taskRun, err := r.framework.TknClient.TaskRuns(r.taskRun.Namespace).Create(context.Background(), r.taskRun, v1.CreateOptions{})
	Expect(err).ShouldNot(HaveOccurred())
	r.taskRun = taskRun
	r.framework.ManageTaskRuns(r.taskRun)
	return r
}

func (r *TaskRunRunner) ExpectFailure() *TaskRunRunner {
	r.taskRun, r.logs = tekton.WaitForTaskRunState(r.framework.Clients, r.taskRun.Namespace, r.taskRun.Name,
		r.taskRun.GetTimeout(context.Background())+constants.Timeouts.TaskRunExtraWaitDelay.Duration,
		tkntest.TaskRunFailed(r.taskRun.Name))
	return r
}

func (r *TaskRunRunner) WaitForTaskRunFinish() *TaskRunRunner {
	r.taskRun, r.logs = tekton.WaitForTaskRunState(r.framework.Clients, r.taskRun.Namespace, r.taskRun.Name,
		r.taskRun.GetTimeout(context.Background())+constants.Timeouts.TaskRunExtraWaitDelay.Duration,
		func(accessor apis.ConditionAccessor) (bool, error) {
			succeeded, _ := tkntest.TaskRunSucceed(r.taskRun.Name)(accessor)
			return succeeded, nil
		})
	return r
}

func (r *TaskRunRunner) ExpectSuccess() *TaskRunRunner {
	r.taskRun, r.logs = tekton.WaitForTaskRunState(r.framework.Clients, r.taskRun.Namespace, r.taskRun.Name,
		r.taskRun.GetTimeout(context.Background())+constants.Timeouts.TaskRunExtraWaitDelay.Duration,
		tkntest.TaskRunSucceed(r.taskRun.Name))
	return r
}

func (r *TaskRunRunner) ExpectSuccessOrFailure(expectSuccess bool) *TaskRunRunner {
	if expectSuccess {
		r.ExpectSuccess()
	} else {
		r.ExpectFailure()
	}
	return r
}

func (r *TaskRunRunner) ExpectLogs(logs ...string) *TaskRunRunner {
	if len(logs) > 0 {
		for _, snippet := range logs {
			Expect(r.logs).Should(ContainSubstring(snippet))
		}
	}
	return r
}

func (r *TaskRunRunner) ExpectTermination(termination *testconfigs.TaskRunExpectedTermination) *TaskRunRunner {
	if termination != nil {
		Expect(r.taskRun.Status.Steps[0].Terminated.ExitCode).Should(Equal(termination.ExitCode))
	}

	return r
}

func (r *TaskRunRunner) GetResults() map[string]string {
	return tekton.TaskResultsToMap(r.taskRun.Status.Results)
}

func (r *TaskRunRunner) ExpectResults(results map[string]string) *TaskRunRunner {
	return r.ExpectResultsWithLen(results, len(results))
}

func (r *TaskRunRunner) ExpectResultsWithLen(results map[string]string, expectedLen int) *TaskRunRunner {
	receivedResults := r.GetResults()

	Expect(receivedResults).Should(HaveLen(expectedLen))

	for resultKey, resultValue := range results {
		Expect(receivedResults[resultKey]).To(Equal(resultValue))
	}
	return r
}
