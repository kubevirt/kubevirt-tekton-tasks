package tekton

import (
	"bytes"
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/onsi/gomega"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1beta1"
	tkntest "github.com/tektoncd/pipeline/test"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/yaml"
	"time"
)

func WaitForTaskRunState(client pipev1beta1.TektonV1beta1Interface, namespace, name string, timeout time.Duration, inState tkntest.ConditionAccessorFn) *v1beta1.TaskRun {
	err := wait.PollImmediate(constants.PollInterval, timeout, func() (bool, error) {
		r, err := client.TaskRuns(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return inState(&r.Status)
	})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	taskRun, err := client.TaskRuns(namespace).Get(name, metav1.GetOptions{})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	return taskRun
}

func GetTaskRunLogs(coreClient clientv1.CoreV1Interface, taskRun *v1beta1.TaskRun) string {
	if taskRun.Status.PodName == "" {
		return ""
	}

	// print logs
	req := coreClient.Pods(taskRun.Namespace).GetLogs(taskRun.Status.PodName, &v1.PodLogOptions{})

	podLogs, err := req.Stream()
	if err != nil {
		return ""
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return ""
	}
	return buf.String()
}

func PrintTaskRunDebugInfo(tknClient pipev1beta1.TektonV1beta1Interface, coreClient clientv1.CoreV1Interface, taskRunNamespace, taskRunName string) {
	// print conditions
	taskRun, err := tknClient.TaskRuns(taskRunNamespace).Get(taskRunName, metav1.GetOptions{})
	if err == nil {
		conditions, _ := yaml.Marshal(taskRun.Status.Conditions)
		fmt.Printf("taskrun conditions:\n%v\n", string(conditions))

		if taskRun.Status.PodName == "" {
			return
		}
		fmt.Printf("%v pod logs:\n%v\n", taskRun.Status.PodName, GetTaskRunLogs(coreClient, taskRun))
	}
}
