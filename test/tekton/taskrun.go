package tekton

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework/clients"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	pipev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tkntest "github.com/tektoncd/pipeline/test"
	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/yaml"
)

func WaitForTaskRunState(clients *clients.Clients, namespace, name string, timeout time.Duration, inState tkntest.ConditionAccessorFn) (*pipev1.TaskRun, string) {
	isCapturing := false
	logs := make(chan string, 1)
	var taskRun *pipev1.TaskRun
	err := wait.PollImmediate(constants.PollInterval, timeout, func() (bool, error) {
		var err error
		taskRun, err = clients.TknClient.TaskRuns(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}

		if taskRun.Status.PodName != "" && !isCapturing {
			req := clients.CoreV1Client.Pods(taskRun.Namespace).GetLogs(taskRun.Status.PodName, &v1.PodLogOptions{
				Follow: true,
			})

			podLogs, err := req.Stream(context.Background())
			if err == nil {
				isCapturing = true
				go func() {
					defer podLogs.Close()
					defer GinkgoRecover()

					result, err := ioutil.ReadAll(podLogs)
					logs <- string(result)
					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				}()
			}
		}
		return inState(&taskRun.Status)
	})
	if err != nil {
		fmt.Printf("%#v \n", taskRun)
	}
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	if isCapturing {
		return taskRun, <-logs
	}

	return taskRun, ""
}

func PrintTaskRunDebugInfo(clients *clients.Clients, taskRunNamespace, taskRunName string) {
	// print conditions
	taskRun, err := clients.TknClient.TaskRuns(taskRunNamespace).Get(context.Background(), taskRunName, metav1.GetOptions{})
	if err == nil {
		conditions, _ := yaml.Marshal(taskRun.Status.Conditions)
		fmt.Printf("taskrun conditions:\n%v\n", string(conditions))

		if taskRun.Status.PodName == "" {
			return
		}
		fmt.Printf("%v pod logs:\n%v\n", taskRun.Status.PodName, getTaskRunLogs(clients.CoreV1Client, taskRun))
	}
}

func getTaskRunLogs(coreClient clientv1.CoreV1Interface, taskRun *pipev1.TaskRun) string {
	if taskRun.Status.PodName == "" {
		return ""
	}

	// print logs
	req := coreClient.Pods(taskRun.Namespace).GetLogs(taskRun.Status.PodName, &v1.PodLogOptions{})

	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return ""
	}
	defer podLogs.Close()

	result, err := ioutil.ReadAll(podLogs)
	if err != nil {
		return ""
	}
	return string(result)
}
