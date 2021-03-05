package runner

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	framework2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/tekton"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/gomega"
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tkntest "github.com/tektoncd/pipeline/test"
	"knative.dev/pkg/apis"
)

type TaskRunRunner struct {
	framework *framework2.Framework
	taskRun   *pipev1beta1.TaskRun
}

func NewTaskRunRunner(framework *framework2.Framework, taskRun *pipev1beta1.TaskRun) *TaskRunRunner {
	Expect(taskRun).ShouldNot(BeNil())
	return &TaskRunRunner{
		framework: framework,
		taskRun:   taskRun,
	}
}

func (r *TaskRunRunner) GetTaskRun() *pipev1beta1.TaskRun {
	return r.taskRun
}

func (r *TaskRunRunner) CreateTaskRun() *TaskRunRunner {
	taskRun, err := r.framework.TknClient.TaskRuns(r.taskRun.Namespace).Create(r.taskRun)
	Expect(err).ShouldNot(HaveOccurred())
	r.taskRun = taskRun
	r.framework.ManageTaskRun(r.taskRun)
	return r
}

func (r *TaskRunRunner) ExpectFailure() *TaskRunRunner {
	r.taskRun = tekton.WaitForTaskRunState(r.framework.TknClient, r.taskRun.Namespace, r.taskRun.Name,
		r.taskRun.GetTimeout()+constants.Timeouts.TaskRunExtraWaitDelay.Duration,
		tkntest.TaskRunFailed(r.taskRun.Name))
	return r
}

func (r *TaskRunRunner) WaitForTaskRunFinish() *TaskRunRunner {
	r.taskRun = tekton.WaitForTaskRunState(r.framework.TknClient, r.taskRun.Namespace, r.taskRun.Name,
		r.taskRun.GetTimeout()+constants.Timeouts.TaskRunExtraWaitDelay.Duration,
		func(accessor apis.ConditionAccessor) (bool, error) {
			succeeded, _ := tkntest.TaskRunSucceed(r.taskRun.Name)(accessor)
			return succeeded, nil
		})
	return r
}

func (r *TaskRunRunner) ExpectSuccess() *TaskRunRunner {
	r.taskRun = tekton.WaitForTaskRunState(r.framework.TknClient, r.taskRun.Namespace, r.taskRun.Name,
		r.taskRun.GetTimeout()+constants.Timeouts.TaskRunExtraWaitDelay.Duration,
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
		taskRunLogs := tekton.GetTaskRunLogs(r.framework.CoreV1Client, r.taskRun)
		for _, snippet := range logs {
			Expect(taskRunLogs).Should(ContainSubstring(snippet))
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
	return tekton.TaskResultsToMap(r.taskRun.Status.TaskRunResults)
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
