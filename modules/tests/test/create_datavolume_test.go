package test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/dv"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/tekton"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/utils"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1beta1"
	tkntest "github.com/tektoncd/pipeline/test"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	cdiv1beta12 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
	cdicliv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"
)

var _ = Describe("Create DataVolume", func() {
	var testConfig *utils.TestConfig
	var taskRunClient pipev1beta1.TaskRunInterface
	var podClient clientv1.PodInterface
	var cdiClientSet cdicliv1beta1.CdiV1beta1Interface

	BeforeEach(func() {
		var err error
		testConfig, err = utils.Setup()
		Expect(err).ShouldNot(HaveOccurred())
		tknClientset, err := versioned.NewForConfig(testConfig.RestConfig)
		Expect(err).ShouldNot(HaveOccurred())
		taskRunClient = tknClientset.TektonV1beta1().TaskRuns(testConfig.DeployNamespace)

		cdiClientSet, err = cdicliv1beta1.NewForConfig(testConfig.RestConfig)
		Expect(err).ShouldNot(HaveOccurred())

		kubeClient, err := clientv1.NewForConfig(testConfig.RestConfig)
		Expect(err).ShouldNot(HaveOccurred())
		podClient = kubeClient.Pods(testConfig.DeployNamespace)
	})

	table.DescribeTable("taskrun fails and no TestDataVolume is created", func(config *dv.CreateDVTestConfig) {
		testConfig.LimitScope(config.LimitScope)
		taskRun, err := config.Init(testConfig).AsTaskRun()
		Expect(err).ShouldNot(HaveOccurred())

		taskRun, taskRunName := tekton.CreateTaskRun(taskRunClient, taskRun)
		defer tekton.DeleteTaskRun(taskRunClient, podClient, taskRunName, testConfig.Debug)

		taskRun = tekton.WaitForTaskRunState(taskRunClient, taskRunName, config.GetTaskRunTimeout().Duration,
			tkntest.TaskRunFailed(taskRunName))

		Expect(taskRun.Status.TaskRunResults).To(BeEmpty())
		if config.ExpectedLogs != "" {
			Expect(tekton.GetTaskRunLogs(podClient, taskRun)).Should(ContainSubstring(config.ExpectedLogs))
		}

		dv := config.TaskData.Datavolume

		if dv != nil && dv.Data.Name != "" {
			// test TestDataVolume should not exist - check just to be sure
			_, err := cdiClientSet.DataVolumes(dv.Data.Namespace).Get(dv.Data.Name, metav1.GetOptions{})
			Expect(err).Should(HaveOccurred())
		}
	},
		table.Entry("empty dv", &dv.CreateDVTestConfig{
			TaskRunTestConfig: utils.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeServiceAccountName,
				ExpectedLogs:   "manifest does not contain DataVolume kind",
			},
			TaskData: dv.CreateDVTaskData{
				Datavolume: nil,
			},
		}),
		table.Entry("malformed dv", &dv.CreateDVTestConfig{
			TaskRunTestConfig: utils.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeServiceAccountName,
				ExpectedLogs:   "manifest does not contain DataVolume kind",
			},
			TaskData: dv.CreateDVTaskData{
				Datavolume: dv.NewBlankDataVolume("malformed").WithoutTypeMeta(),
			},
		}),
		table.Entry("no service account", &dv.CreateDVTestConfig{
			TaskRunTestConfig: utils.TaskRunTestConfig{
				ExpectedLogs: "datavolumes.cdi.kubevirt.io is forbidden",
			},
			TaskData: dv.CreateDVTaskData{
				Datavolume: dv.NewBlankDataVolume("no-sa"),
			},
		}),
		table.Entry("missing name", &dv.CreateDVTestConfig{
			TaskRunTestConfig: utils.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeServiceAccountName,
				ExpectedLogs:   "invalid: metadata.name: Required value: name",
			},
			TaskData: dv.CreateDVTaskData{
				Datavolume: dv.NewBlankDataVolume(""),
			},
		}),
		table.Entry("cannot create a TestDataVolume in different namespace (namespace scoped)", &dv.CreateDVTestConfig{
			TaskRunTestConfig: utils.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeServiceAccountName,
				Namespace:      CustomTargetNS,
				LimitScope:     utils.NamespaceScope,
				ExpectedLogs:   "datavolumes.cdi.kubevirt.io is forbidden",
			},
			TaskData: dv.CreateDVTaskData{
				Datavolume: dv.NewBlankDataVolume("namespace-scope"),
			},
		}),
	)

	table.DescribeTable("TestDataVolume and PVC is created successfully", func(config *dv.CreateDVTestConfig) {
		testConfig.LimitScope(config.LimitScope)
		taskRun, err := config.Init(testConfig).AsTaskRun()
		Expect(err).ShouldNot(HaveOccurred())

		taskRun, taskRunName := tekton.CreateTaskRun(taskRunClient, taskRun)
		defer tekton.DeleteTaskRun(taskRunClient, podClient, taskRunName, testConfig.Debug)

		taskRun = tekton.WaitForTaskRunState(taskRunClient, taskRunName, config.GetTaskRunTimeout().Duration,
			tkntest.TaskRunSucceed(taskRunName))

		results := tekton.TaskResultsToMap(taskRun.Status.TaskRunResults)

		Expect(results).Should(HaveLen(2))
		dvName := results[CreateVMFromManifestResults.Name]
		dvNamespace := results[CreateVMFromManifestResults.Namespace]
		Expect(dvName).ToNot(BeEmpty())
		Expect(dvNamespace).ToNot(BeEmpty())
		defer dv.DeleteDataVolume(cdiClientSet.DataVolumes(dvNamespace), dvName, testConfig.Debug)

		var dataVolume *cdiv1beta12.DataVolume
		timeout := config.GetWaitForDVTimeout()
		Expect(timeout).ToNot(BeNil())

		err = wait.PollImmediate(PollInterval, timeout.Duration, func() (bool, error) {
			dataVolume, err = cdiClientSet.DataVolumes(dvNamespace).Get(dvName, metav1.GetOptions{})
			if err != nil {
				return true, err
			}
			return dv.GetConditionMap(dataVolume)[cdiv1beta12.DataVolumeBound] == v1.ConditionTrue &&
				dataVolume.Status.Phase == cdiv1beta12.Succeeded, nil
		})
		Expect(err).ShouldNot(HaveOccurred())
		if config.ExpectedLogs != "" {
			Expect(tekton.GetTaskRunLogs(podClient, taskRun)).Should(ContainSubstring(config.ExpectedLogs))
		}
	},
		table.Entry("blank wait", &dv.CreateDVTestConfig{
			TaskRunTestConfig: utils.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeServiceAccountName,
				Timeout:        Timeouts.SmallBlankDVCreation,
				ExpectedLogs:   "Created",
			},
			TaskData: dv.CreateDVTaskData{
				Datavolume:     dv.NewBlankDataVolume("blank"),
				WaitForSuccess: true,
			},
		}),
		table.Entry("blank no wait", &dv.CreateDVTestConfig{
			TaskRunTestConfig: utils.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeServiceAccountName,
				Timeout:        Timeouts.SmallBlankDVCreation,
				ExpectedLogs:   "Created",
			},
			TaskData: dv.CreateDVTaskData{
				Datavolume:     dv.NewBlankDataVolume("blank-wait"),
				WaitForSuccess: true,
			},
		}),
		table.Entry("works also in the same namespace as deploy (cluster scoped)", &dv.CreateDVTestConfig{
			TaskRunTestConfig: utils.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeServiceAccountName,
				Namespace:      DeployTargetNS,
				LimitScope:     utils.ClusterScope,
				Timeout:        Timeouts.SmallBlankDVCreation,
				ExpectedLogs:   "Created",
			},
			TaskData: dv.CreateDVTaskData{
				Datavolume:     dv.NewBlankDataVolume("namespace-scope"),
				WaitForSuccess: true,
			},
		}),
	)
})
