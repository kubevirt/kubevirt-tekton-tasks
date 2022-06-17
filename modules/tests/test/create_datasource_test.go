package test

import (
	"context"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datavolume"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/ds"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/dv"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Create DataSource", func() {
	f := framework.NewFramework()
	DescribeTable("DataSource is created successfully", func(config *testconfigs.CreateDVTestConfig) {
		f.TestSetup(config)
		if config.DataSource.TestDV != nil {
			newDataVolume, err := f.CdiClient.DataVolumes(config.DataSource.TestDV.Namespace).Create(context.TODO(), config.DataSource.TestDV, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataVolumes(newDataVolume)

			err = dv.WaitForSuccessfulDataVolume(f.CdiClient, newDataVolume.Namespace, newDataVolume.Name, config.GetWaitForDVTimeout())
			Expect(err).ShouldNot(HaveOccurred())
		}

		datasource := config.DataSource.DataSource
		f.ManageDataSources(datasource)

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResults(map[string]string{
				CreateDataVolumeFromManifestResults.Name:      datasource.Name,
				CreateDataVolumeFromManifestResults.Namespace: datasource.Namespace,
			})

		err := ds.WaitForSuccessfulDataSource(f.CdiClient, datasource.Namespace, datasource.Name, config.GetWaitForDSTimeout())
		Expect(err).ShouldNot(HaveOccurred())
	},
		Entry("pointing to pvc with wait", &testconfigs.CreateDVTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeFromManifestServiceAccountName,
				Timeout:        Timeouts.SmallDVCreation,
				ExpectedLogs:   "Created",
			},
			DataSource: testconfigs.CreateDSTaskData{
				DataSource:     ds.NewDataSource("test-ds"),
				WaitForSuccess: true,
				Namespace:      DeployTargetNS,
				TestDV:         datavolume.NewBlankDataVolume("test-ds").Build(),
			},
		}),
		Entry("pointing to pvc without wait", &testconfigs.CreateDVTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateDataVolumeFromManifestServiceAccountName,
				Timeout:        Timeouts.SmallDVCreation,
				ExpectedLogs:   "Created",
			},
			DataSource: testconfigs.CreateDSTaskData{
				DataSource:     ds.NewDataSource("test-ds"),
				WaitForSuccess: false,
				Namespace:      DeployTargetNS,
				TestDV:         datavolume.NewBlankDataVolume("test-ds").Build(),
			},
		}),
	)
})
