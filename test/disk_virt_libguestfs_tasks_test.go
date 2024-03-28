package test

import (
	"context"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datavolume"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/dataobject"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	customizeCommand = `run-command date
mkdir /test
touch /test/whale
append-line /test/whale:hello blue world
run-command date
no-logfile
`
)

var _ = Describe("Run disk virt-customize / virt-sysprep", func() {
	f := framework.NewFramework()
	for _, t := range []constants.LibguestfsTaskType{constants.VirtCustomizeTaskType, constants.VirtSysPrepTaskType} {
		taskType := t
		DescribeTable(string(taskType)+" taskrun fails", func(config *testconfigs.DiskVirtLibguestfsTestConfig) {
			config.TaskData.LibguestfsTaskType = taskType
			f.TestSetup(config)

			if dataVolume := config.TaskData.Datavolume; dataVolume != nil {
				dataVolume, err := f.CdiClient.DataVolumes(dataVolume.Namespace).Create(context.Background(), dataVolume, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageDataVolumes(dataVolume)
				err = dataobject.WaitForSuccessfulDataVolume(f.KubevirtClient, dataVolume.Namespace, dataVolume.Name, constants.Timeouts.SmallDVCreation.Duration)
				Expect(err).ShouldNot(HaveOccurred())
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(nil)
		},
			Entry("no pvc", &testconfigs.DiskVirtLibguestfsTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
				TaskData:          testconfigs.DiskVirtLibguestfsTaskData{},
			}),
			Entry("invalid pvc", &testconfigs.DiskVirtLibguestfsTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					Timeout: &metav1.Duration{time.Minute},
				},
				TaskData: testconfigs.DiskVirtLibguestfsTaskData{
					PVCName: "illegal-test-pvc-1548748",
				},
			}),
			Entry("no commands", &testconfigs.DiskVirtLibguestfsTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "env variable is required",
				},
				TaskData: testconfigs.DiskVirtLibguestfsTaskData{
					Datavolume: datavolume.NewBlankDataVolume("no-customize-commands").Build(),
				},
			}),
			Entry("wrong customize commands", &testconfigs.DiskVirtLibguestfsTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "command 'illegal-operation' not valid",
				},
				TaskData: testconfigs.DiskVirtLibguestfsTaskData{
					Datavolume: datavolume.NewBlankDataVolume("wrong-customize-commands").Build(),
					Commands:   "illegal-operation",
				},
			}),
			Entry("wrong additional options", &testconfigs.DiskVirtLibguestfsTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: "unrecognized option '--illegal-command'",
				},
				TaskData: testconfigs.DiskVirtLibguestfsTaskData{
					Datavolume:        datavolume.NewBlankDataVolume("wrong-additional-options").Build(),
					Commands:          "update",
					AdditionalOptions: "--verbose --illegal-command illegal args",
				},
			}),
			Entry("empty disk fails", &testconfigs.DiskVirtLibguestfsTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogsList: []string{
						"no operating systems were found in the guest image",
					},
				},
				TaskData: testconfigs.DiskVirtLibguestfsTaskData{
					Datavolume: datavolume.NewBlankDataVolume("empty-disk").Build(),
					Commands:   "update",
				},
			}),
		)

		It(string(taskType)+" works", func() {
			// main test
			testConfig := &testconfigs.DiskVirtLibguestfsTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogsList: []string{
						"Running: date",
						"Making directory: /test",
						"Running touch: /test/whale",
						"Appending line to /test/whale",
						"Scrubbing the log file",
					},
				},
				TaskData: testconfigs.DiskVirtLibguestfsTaskData{
					Datavolume: datavolume.NewBlankDataVolume("basic-functionality").
						WithSize(5, resource.Giga).
						WithRegistrySource("docker://kubevirt/fedora-cloud-registry-disk-demo").
						Build(),
					Commands:           customizeCommand,
					LibguestfsTaskType: taskType,
				},
			}
			f.TestSetup(testConfig)

			dataVolume := testConfig.TaskData.Datavolume
			// prepare DataVolume
			dataVolume, err := f.CdiClient.DataVolumes(dataVolume.Namespace).Create(context.Background(), dataVolume, metav1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataVolumes(dataVolume)
			err = dataobject.WaitForSuccessfulDataVolume(f.KubevirtClient, dataVolume.Namespace, dataVolume.Name, constants.Timeouts.SmallDVCreation.Duration)
			Expect(err).ShouldNot(HaveOccurred())

			// run main test
			runner.NewTaskRunRunner(f, testConfig.GetTaskRunWithName("-test")).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(testConfig.GetAllExpectedLogs()...).
				ExpectResults(nil)

			// test the result data is preserved on the PVC

			expectConfig := &testconfigs.DiskVirtLibguestfsTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogsList: []string{
						"\"level\":\"debug\"",
						"-x",
						"Running: cat /test/whale",
						"hello blue world",
					},
				},
				TaskData: testconfigs.DiskVirtLibguestfsTaskData{
					PVCName:            dataVolume.Name,
					Commands:           "run-command cat /test/whale",
					Verbose:            true,
					AdditionalOptions:  "--no-network --verbose",
					LibguestfsTaskType: taskType,
				},
			}

			f.TestSetup(expectConfig)

			runner.NewTaskRunRunner(f, expectConfig.GetTaskRunWithName("-expect")).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(expectConfig.GetAllExpectedLogs()...).
				ExpectResults(nil)
		})
	}
})
