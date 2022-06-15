package test

import (
	"context"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datavolume"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/dv"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Create DataVolume", func() {
	f := framework.NewFramework()

	DescribeTable("taskrun fails and no DataVolume is created", func(config *testconfigs.CreateDVTestConfig) {
		f.TestSetup(config)

		dataVolume := config.TaskData.Datavolume
		f.ManageDataVolumes(dataVolume)

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectFailure().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResults(nil)

		if dataVolume != nil && dataVolume.Name != "" && dataVolume.Namespace != "" {
			_, err := f.CdiClient.DataVolumes(dataVolume.Namespace).Get(context.TODO(), dataVolume.Name, metav1.GetOptions{})
			Expect(err).Should(HaveOccurred())
		}
	},
		Entry("empty dv", &testconfigs.CreateDVTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeFromManifestServiceAccountName,
				ExpectedLogs:   "manifest does not contain DataVolume or DataSource kind",
			},
			TaskData: testconfigs.CreateDVTaskData{
				Datavolume: nil,
			},
		}),
		Entry("malformed dv", &testconfigs.CreateDVTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeFromManifestServiceAccountName,
				ExpectedLogs:   "manifest does not contain DataVolume or DataSource kind",
			},
			TaskData: testconfigs.CreateDVTaskData{
				Datavolume: datavolume.NewBlankDataVolume("malformed").WithoutTypeMeta().Build(),
			},
		}),
		Entry("missing name", &testconfigs.CreateDVTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeFromManifestServiceAccountName,
				ExpectedLogs:   "invalid: metadata.name: Required value: name",
			},
			TaskData: testconfigs.CreateDVTaskData{
				Datavolume: datavolume.NewBlankDataVolume("").Build(),
			},
		}),
		Entry("[NAMESPACE SCOPED] cannot create a DataVolume in different namespace", &testconfigs.CreateDVTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeFromManifestServiceAccountName,
				LimitTestScope: NamespaceTestScope,
				ExpectedLogs:   "datavolumes.cdi.kubevirt.io is forbidden",
			},
			TaskData: testconfigs.CreateDVTaskData{
				Datavolume: datavolume.NewBlankDataVolume("different-ns-namespace-scope").Build(),
				Namespace:  SystemTargetNS,
			},
		}),
	)

	DescribeTable("DataVolume and PVC is created successfully", func(config *testconfigs.CreateDVTestConfig) {
		f.TestSetup(config)

		dataVolume := config.TaskData.Datavolume
		f.ManageDataVolumes(dataVolume)

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResults(map[string]string{
				CreateDataVolumeFromManifestResults.Name:      dataVolume.Name,
				CreateDataVolumeFromManifestResults.Namespace: dataVolume.Namespace,
			})

		err := dv.WaitForSuccessfulDataVolume(f.CdiClient, dataVolume.Namespace, dataVolume.Name, config.GetWaitForDVTimeout())
		Expect(err).ShouldNot(HaveOccurred())
	},
		Entry("blank wait", &testconfigs.CreateDVTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeFromManifestServiceAccountName,
				Timeout:        Timeouts.SmallDVCreation,
				ExpectedLogs:   "Created",
			},
			TaskData: testconfigs.CreateDVTaskData{
				Datavolume:     datavolume.NewBlankDataVolume("blank").Build(),
				WaitForSuccess: true,
			},
		}),
		Entry("blank no wait", &testconfigs.CreateDVTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeFromManifestServiceAccountName,
				Timeout:        Timeouts.SmallDVCreation,
				ExpectedLogs:   "Created",
			},
			TaskData: testconfigs.CreateDVTaskData{
				Datavolume:     datavolume.NewBlankDataVolume("blank-wait").Build(),
				WaitForSuccess: true,
			},
		}),
		Entry("[CLUSTER SCOPED] works also in the same namespace as deploy", &testconfigs.CreateDVTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeFromManifestServiceAccountName,
				LimitTestScope: ClusterTestScope,
				Timeout:        Timeouts.SmallDVCreation,
				ExpectedLogs:   "Created",
			},
			TaskData: testconfigs.CreateDVTaskData{
				Datavolume:     datavolume.NewBlankDataVolume("same-ns-cluster-scope").Build(),
				WaitForSuccess: true,
				Namespace:      DeployTargetNS,
			},
		}),
	)

	It("taskrun fails and DataVolume is created but does not import successfully", func() {
		config := &testconfigs.CreateDVTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeFromManifestServiceAccountName,
				Timeout:        Timeouts.QuickTaskRun,
			},
			TaskData: testconfigs.CreateDVTaskData{
				Datavolume: datavolume.NewBlankDataVolume("blank").
					WithURLSource("https://invalid.source.my.domain.fail").Build(),
				WaitForSuccess: true,
			},
		}
		f.TestSetup(config)

		dataVolume := config.TaskData.Datavolume
		f.ManageDataVolumes(dataVolume)

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccess()

		Eventually(func() bool {
			d, err := f.CdiClient.DataVolumes(dataVolume.Namespace).Get(context.TODO(), dataVolume.Name, metav1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(d.Spec.Source.HTTP.URL).To(Equal(dataVolume.Spec.Source.HTTP.URL))

			return dv.HasDataVolumeFailedToImport(d)
		}, 30*time.Second, 1*time.Second).Should(BeTrue())

	})
})
