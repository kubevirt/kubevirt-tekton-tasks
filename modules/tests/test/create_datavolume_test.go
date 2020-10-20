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

	table.DescribeTable("taskrun fails and no DV is created", func(config *dv.CreateDVTestConfig) {
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

		if config.Datavolume != nil && config.Datavolume.Name != "" {
			// test DV should not exist - check just to be sure
			_, err := cdiClientSet.DataVolumes(config.Datavolume.Namespace).Get(config.Datavolume.Name, metav1.GetOptions{})
			Expect(err).Should(HaveOccurred())
		}
	},
		table.Entry("empty dv", &dv.CreateDVTestConfig{
			Datavolume:     nil,
			ServiceAccount: CreateDataVolumeServiceAccountName,
			ExpectedLogs:   "manifest does not contain DataVolume kind",
		}),
		table.Entry("malformed dv", &dv.CreateDVTestConfig{
			Datavolume:     dv.NewBlankDV("malformed").WithoutTypeMeta().Build(),
			ServiceAccount: CreateDataVolumeServiceAccountName,
			ExpectedLogs:   "manifest does not contain DataVolume kind",
		}),
		table.Entry("no service account", &dv.CreateDVTestConfig{
			Datavolume:   dv.NewBlankDV("no-sc").Build(),
			ExpectedLogs: "datavolumes.cdi.kubevirt.io is forbidden",
		}),
		table.Entry("missing name", &dv.CreateDVTestConfig{
			Datavolume:     dv.NewBlankDV("").Build(),
			ServiceAccount: CreateDataVolumeServiceAccountName,
			ExpectedLogs:   "invalid: metadata.name: Required value: name",
		}),
		table.Entry("cannot create a DV in different namespace (namespace scoped)", &dv.CreateDVTestConfig{
			Datavolume:     dv.NewBlankDV("namespace-scope").Build(),
			ServiceAccount: CreateDataVolumeServiceAccountName,
			Namespace:      CustomTargetNS,
			LimitScope:     utils.NamespaceScope,
			ExpectedLogs:   "datavolumes.cdi.kubevirt.io is forbidden",
		}),
	)

	table.DescribeTable("DV and PVC is created successfully", func(config *dv.CreateDVTestConfig) {
		testConfig.LimitScope(config.LimitScope)
		taskRun, err := config.Init(testConfig).AsTaskRun()
		Expect(err).ShouldNot(HaveOccurred())

		taskRun, taskRunName := tekton.CreateTaskRun(taskRunClient, taskRun)
		defer tekton.DeleteTaskRun(taskRunClient, podClient, taskRunName, testConfig.Debug)

		taskRun = tekton.WaitForTaskRunState(taskRunClient, taskRunName, config.GetTaskRunTimeout().Duration,
			tkntest.TaskRunSucceed(taskRunName))

		results := tekton.TaskResultsToMap(taskRun.Status.TaskRunResults)

		Expect(results).Should(HaveLen(2))
		dvName := results[CreateDataVolumeFromManifestResults.Name]
		dvNamespace := results[CreateDataVolumeFromManifestResults.Namespace]
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
			Datavolume:     dv.NewBlankDV("blank").Build(),
			ServiceAccount: CreateDataVolumeServiceAccountName,
			Timeout:        Timeouts.SmallBlankDVCreation,
			WaitForSuccess: true,
			ExpectedLogs:   "Created",
		}),
		table.Entry("blank no wait", &dv.CreateDVTestConfig{
			Datavolume:     dv.NewBlankDV("blank-wait").Build(),
			ServiceAccount: CreateDataVolumeServiceAccountName,
			Timeout:        Timeouts.SmallBlankDVCreation,
			ExpectedLogs:   "Created",
		}),
		table.Entry("works also in the same namespace as deploy (cluster scoped)", &dv.CreateDVTestConfig{
			Datavolume:     dv.NewBlankDV("namespace-scope").Build(),
			ServiceAccount: CreateDataVolumeServiceAccountName,
			Namespace:      DeployTargetNS,
			LimitScope:     utils.ClusterScope,
			Timeout:        Timeouts.SmallBlankDVCreation,
			ExpectedLogs:   "Created",
		}),
	)
})
