package test

import (
	"context"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	testtemplate "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	templatev1 "github.com/openshift/api/template/v1"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

type ExpectedResults struct {
	ExpectedCPUSocketsTopologyNumber uint32
	ExpectedCPUCoresTopologyNumber   uint32
	ExpectedCPUThreadsTopologyNumber uint32
	ExpectedTemplateLabels           map[string]string
	ExpectedTemplateAnnotations      map[string]string
	ExpectedVMLabels                 map[string]string
	ExpectedVMAnnotations            map[string]string
	ExpectedVMMemory                 resource.Quantity
	ExpectedDisks                    []kubevirtv1.Disk
	ExpectedVolumes                  []kubevirtv1.Volume
	ExpectedDataVolumeTemplates      []kubevirtv1.DataVolumeTemplateSpec
}

func (e ExpectedResults) ExpectCPUTopology(vm *kubevirtv1.VirtualMachine) {
	Expect(vm.Spec.Template.Spec.Domain.CPU.Sockets).To(Equal(e.ExpectedCPUSocketsTopologyNumber), "cpu sockets should equal")
	Expect(vm.Spec.Template.Spec.Domain.CPU.Cores).To(Equal(e.ExpectedCPUCoresTopologyNumber), "cpu cores should equal")
	Expect(vm.Spec.Template.Spec.Domain.CPU.Threads).To(Equal(e.ExpectedCPUThreadsTopologyNumber), "cpu threads should equal")
}

func (e ExpectedResults) ExpectTemplateLabels(template *templatev1.Template) {
	checkValuesInMap(template.Labels, e.ExpectedTemplateLabels, "template labels")
}

func (e ExpectedResults) ExpectTemplateAnnotations(template *templatev1.Template) {
	checkValuesInMap(template.Annotations, e.ExpectedTemplateAnnotations, "template annotations")
}

func (e ExpectedResults) ExpectVMLabels(vm *kubevirtv1.VirtualMachine) {
	checkValuesInMap(vm.Labels, e.ExpectedVMLabels, "vm labels")
}

func (e ExpectedResults) ExpectVMAnnotations(vm *kubevirtv1.VirtualMachine) {
	checkValuesInMap(vm.Annotations, e.ExpectedVMAnnotations, "vm annotations")
}

func (e ExpectedResults) ExpectVMMemory(vm *kubevirtv1.VirtualMachine) {
	Expect(vm.Spec.Template.Spec.Domain.Resources.Requests.Memory()).To(Equal(&e.ExpectedVMMemory), "memory should equal")
}

func (e ExpectedResults) ExpectVMDisks(vm *kubevirtv1.VirtualMachine) {
	Expect(vm.Spec.Template.Spec.Domain.Devices.Disks).To(Equal(e.ExpectedDisks), "vm disks should equal")
}

func (e ExpectedResults) ExpectVMVolumes(vm *kubevirtv1.VirtualMachine) {
	Expect(vm.Spec.Template.Spec.Volumes).To(Equal(e.ExpectedVolumes), "vm volumes should equal")
}

func (e ExpectedResults) ExpectDataVolumeTemplates(vm *kubevirtv1.VirtualMachine) {
	Expect(vm.Spec.DataVolumeTemplates[0]).To(Equal(e.ExpectedDataVolumeTemplates[0]), "vm datavolume templates should equal")
}

func checkValuesInMap(updatedMap, expectedResult map[string]string, message string) {
	for key, value := range expectedResult {
		Expect(updatedMap[key]).To(Equal(value), "value should equal in "+message)
	}
}

var _ = Describe("Modify template task", func() {
	f := framework.NewFramework().LimitEnvScope(OKDEnvScope)
	Context("modify template fail", func() {
		DescribeTable("taskrun fails and no template is created", func(config *testconfigs.ModifyTemplateTestConfig) {
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				t, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(t)
			}
			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectFailure().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(nil)
		},
			Entry("no source template name specified", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "template-name param has to be specified",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			Entry("no template specified", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "templates.template.openshift.io \"cirros-vm-template\" not found",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
				},
			}),
			Entry("cannot updated template in different namespace", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountNameNamespaced,
					ExpectedLogs:   "cannot get resource \"templates\" in API group \"template.openshift.io\"",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:            testtemplate.CirrosTemplateName,
					SourceTemplateNamespace: SystemTargetNS,
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				},
			}),
			Entry("wrong memory value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "quantities must match the regular expression",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					Memory:       "wrong memory value",
				},
			}),
			Entry("wrong number of CPU cores value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "parsing \"wrong cpu cores\": invalid syntax",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					CPUCores:     "wrong cpu cores",
				},
			}),
			Entry("wrong number of CPU Sockets value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "parsing \"wrong cpu sockets\": invalid syntax",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					CPUCores:     CPUCoresTopologyNumberStr,
					CPUThreads:   CPUThreadsTopologyNumberStr,
					CPUSockets:   "wrong cpu sockets",
				},
			}),
			Entry("wrong number of CPU Threads value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "parsing \"wrong cpu threads\": invalid syntax",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					CPUCores:     CPUCoresTopologyNumberStr,
					CPUThreads:   "wrong cpu threads",
				},
			}),
			Entry("wrong template labels value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "pair should be in \"KEY:VAL\" format",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:   testtemplate.CirrosTemplateName,
					CPUCores:       CPUCoresTopologyNumberStr,
					CPUThreads:     CPUThreadsTopologyNumberStr,
					CPUSockets:     CPUSocketsTopologyNumberStr,
					TemplateLabels: []string{"singleLabel"},
				},
			}),
			Entry("wrong template annotations value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "pair should be in \"KEY:VAL\" format",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:        testtemplate.CirrosTemplateName,
					CPUCores:            CPUCoresTopologyNumberStr,
					CPUThreads:          CPUThreadsTopologyNumberStr,
					CPUSockets:          CPUSocketsTopologyNumberStr,
					TemplateLabels:      MockTemplateLabels,
					TemplateAnnotations: []string{"singleAnnotation"},
				},
			}),
			Entry("wrong vm labels value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "pair should be in \"KEY:VAL\" format",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:        testtemplate.CirrosTemplateName,
					CPUCores:            CPUCoresTopologyNumberStr,
					CPUThreads:          CPUThreadsTopologyNumberStr,
					CPUSockets:          CPUSocketsTopologyNumberStr,
					TemplateLabels:      MockTemplateLabels,
					TemplateAnnotations: MockTemplateAnnotations,
					VMLabels:            []string{"singleLabel"},
				},
			}),
			Entry("wrong template annotations value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "pair should be in \"KEY:VAL\" format",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:        testtemplate.CirrosTemplateName,
					CPUCores:            CPUCoresTopologyNumberStr,
					CPUThreads:          CPUThreadsTopologyNumberStr,
					CPUSockets:          CPUSocketsTopologyNumberStr,
					TemplateLabels:      MockTemplateLabels,
					TemplateAnnotations: MockTemplateAnnotations,
					VMLabels:            MockVMLabels,
					VMAnnotations:       []string{"singleAnnotation"},
				},
			}),
			Entry("wrong disks value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "invalid character 'w' looking for beginning of value",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					Disks:        WrongStrSlice,
				},
			}),
			Entry("wrong volumes value provided", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "invalid character 'w' looking for beginning of value",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName: testtemplate.CirrosTemplateName,
					Volumes:      WrongStrSlice,
				},
			}),
			Entry("cannot delete non-existent template", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
					ExpectedLogs:   "templates.template.openshift.io \"" + testtemplate.CirrosTemplateName + "\" not found",
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					TemplateName:   testtemplate.CirrosTemplateName,
					DeleteTemplate: true,
				},
			}),
		)
	})
	Context("modify template sucess", func() {
		DescribeTable("taskrun succeded and template is updated", func(config *testconfigs.ModifyTemplateTestConfig, expectedResults ExpectedResults) {
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				t, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(t)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.TemplateName,
					"namespace": config.TaskData.TemplateNamespace,
				})

			template, err := f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.Background(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(template).ToNot(BeNil(), "new template should exists")
			f.ManageTemplates(template)

			vm, _, err := zutils.DecodeVM(template)
			Expect(err).ShouldNot(HaveOccurred(), "decode VM")

			expectedResults.ExpectCPUTopology(vm)
			expectedResults.ExpectTemplateLabels(template)
			expectedResults.ExpectTemplateAnnotations(template)
			expectedResults.ExpectVMLabels(vm)
			expectedResults.ExpectVMAnnotations(vm)
			expectedResults.ExpectVMMemory(vm)
			expectedResults.ExpectVMDisks(vm)
			expectedResults.ExpectVMVolumes(vm)
			expectedResults.ExpectDataVolumeTemplates(vm)

		},
			Entry("should update template in the same namespace", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:            testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateName:        testtemplate.CirrosTemplateName,
					CPUCores:            CPUCoresTopologyNumberStr,
					CPUSockets:          CPUSocketsTopologyNumberStr,
					CPUThreads:          CPUThreadsTopologyNumberStr,
					Memory:              MemoryValue,
					TemplateAnnotations: MockTemplateAnnotations,
					TemplateLabels:      MockTemplateLabels,
					VMAnnotations:       MockVMAnnotations,
					VMLabels:            MockVMLabels,
					Disks:               MockDisks,
					Volumes:             MockVolumes,
					DataVolumeTemplates: MockDataVolumeTemplates,
				},
			}, ExpectedResults{
				ExpectedCPUSocketsTopologyNumber: CPUSocketsTopologyNumber,
				ExpectedCPUCoresTopologyNumber:   CPUCoresTopologyNumber,
				ExpectedCPUThreadsTopologyNumber: CPUThreadsTopologyNumber,
				ExpectedVMMemory:                 resource.MustParse(MemoryValue),
				ExpectedTemplateAnnotations:      TemplateAnnotationsMap,
				ExpectedTemplateLabels:           TemplateLabelsMap,
				ExpectedVMAnnotations:            VMAnnotationsMap,
				ExpectedVMLabels:                 VMLabelsMap,
				ExpectedDisks:                    Disks,
				ExpectedVolumes:                  Volumes,
				ExpectedDataVolumeTemplates:      DataVolumeTemplates,
			}),
			Entry("should update template across namespaces", &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateName:            testtemplate.CirrosTemplateName,
					SourceTemplateNamespace: DeployTargetNS,
					CPUCores:                CPUCoresTopologyNumberStr,
					CPUSockets:              CPUSocketsTopologyNumberStr,
					CPUThreads:              CPUThreadsTopologyNumberStr,
					Memory:                  MemoryValue,
					TemplateAnnotations:     MockTemplateAnnotations,
					TemplateLabels:          MockTemplateLabels,
					VMAnnotations:           MockVMAnnotations,
					VMLabels:                MockVMLabels,
					Disks:                   MockDisks,
					Volumes:                 MockVolumes,
					DataVolumeTemplates:     MockDataVolumeTemplates,
				},
			}, ExpectedResults{
				ExpectedCPUSocketsTopologyNumber: CPUSocketsTopologyNumber,
				ExpectedCPUCoresTopologyNumber:   CPUCoresTopologyNumber,
				ExpectedCPUThreadsTopologyNumber: CPUThreadsTopologyNumber,
				ExpectedVMMemory:                 resource.MustParse(MemoryValue),
				ExpectedTemplateAnnotations:      TemplateAnnotationsMap,
				ExpectedTemplateLabels:           TemplateLabelsMap,
				ExpectedVMAnnotations:            VMAnnotationsMap,
				ExpectedVMLabels:                 VMLabelsMap,
				ExpectedDisks:                    Disks,
				ExpectedVolumes:                  Volumes,
				ExpectedDataVolumeTemplates:      DataVolumeTemplates,
			}),
		)

		It("taskrun succeded and template datavolume is removed and volumes, disks are updated", func() {
			config := &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:                 testtemplate.NewRhelDesktopTinyTemplate().Build(),
					TemplateName:             testtemplate.RhelTemplateName,
					SourceTemplateNamespace:  DeployTargetNS,
					DeleteDatavolumeTemplate: true,
				},
			}
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				t, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(t)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.TemplateName,
					"namespace": config.TaskData.TemplateNamespace,
				})

			template, err := f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.Background(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(template).ToNot(BeNil(), "new template should exists")
			f.ManageTemplates(template)

			vm, _, err := zutils.DecodeVM(template)
			Expect(err).ShouldNot(HaveOccurred(), "decode VM")
			Expect(len(vm.Spec.DataVolumeTemplates)).To(Equal(0), "datavolume template should be empty")

			Expect(len(vm.Spec.Template.Spec.Volumes)).To(Equal(1), "there should be only 1 volume")
			emptyVolume := kubevirtv1.Volume{}
			Expect(vm.Spec.Template.Spec.Volumes[0].DataVolume).To(Equal(emptyVolume.DataVolume), "data volume should be nil")

			Expect(len(vm.Spec.Template.Spec.Domain.Devices.Disks)).To(Equal(1), "there should be only 1 disk")
			Expect(vm.Spec.Template.Spec.Domain.Devices.Disks[0].Name).ToNot(Equal("rootdisk"), "disk should not have name rootdisk")
		})

		It("taskrun succeded and template datavolume is removed and replaced by a new one", func() {
			config := &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:                 testtemplate.NewRhelDesktopTinyTemplate().Build(),
					TemplateName:             testtemplate.RhelTemplateName,
					SourceTemplateNamespace:  DeployTargetNS,
					DataVolumeTemplates:      MockDataVolumeTemplates,
					DeleteDatavolumeTemplate: true,
				},
			}
			f.TestSetup(config)

			if template := config.TaskData.Template; template != nil {
				t, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(t)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.TemplateName,
					"namespace": config.TaskData.TemplateNamespace,
				})

			template, err := f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.Background(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(template).ToNot(BeNil(), "new template should exists")
			f.ManageTemplates(template)

			vm, _, err := zutils.DecodeVM(template)
			Expect(err).ShouldNot(HaveOccurred(), "decode VM")
			Expect(len(vm.Spec.DataVolumeTemplates)).To(Equal(1), "there should be only one datavolume template")
			Expect(vm.Spec.DataVolumeTemplates[0]).To(Equal(DataVolumeTemplates[0]), "vm datavolume templates should equal")
		})

		It("taskrun succeded and defaults do not modify template", func() {
			config := &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:                testtemplate.NewRhelDesktopTinyTemplate().Build(),
					TemplateName:            testtemplate.RhelTemplateName,
					SourceTemplateNamespace: DeployTargetNS,
				},
			}
			f.TestSetup(config)

			template := config.TaskData.Template
			if template != nil {
				template, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(template)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.TemplateName,
					"namespace": config.TaskData.TemplateNamespace,
				})

			updatedTemplate, err := f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.Background(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(updatedTemplate).ToNot(BeNil(), "new template should exists")
			f.ManageTemplates(updatedTemplate)

			vm, _, err := zutils.DecodeVM(template)
			Expect(err).ShouldNot(HaveOccurred(), "decode VM")
			updatedVm, _, err := zutils.DecodeVM(updatedTemplate)
			Expect(err).ShouldNot(HaveOccurred(), "decode updated VM")

			Expect(updatedTemplate.Labels).To(Equal(template.Labels), "templateLabels should be unchanged")
			Expect(updatedTemplate.Annotations).To(Equal(template.Annotations), "templateAnnotations should be unchanged")
			Expect(updatedVm.Spec.Template.Spec.Domain.CPU.Sockets).To(Equal(vm.Spec.Template.Spec.Domain.CPU.Sockets), "cpuSockets should be unchanged")
			Expect(updatedVm.Spec.Template.Spec.Domain.CPU.Cores).To(Equal(vm.Spec.Template.Spec.Domain.CPU.Cores), "cpuCores should be unchanged")
			Expect(updatedVm.Spec.Template.Spec.Domain.CPU.Threads).To(Equal(vm.Spec.Template.Spec.Domain.CPU.Threads), "cpuThreads should be unchanged")
			Expect(updatedVm.Labels).To(Equal(vm.Labels), "vmLabels should be unchanged")
			Expect(updatedVm.Annotations).To(Equal(vm.Annotations), "vmAnnotations should be unchanged")
			Expect(updatedVm.Spec.Template.Spec.Domain.Devices.Disks).To(Equal(vm.Spec.Template.Spec.Domain.Devices.Disks), "disks should be unchanged")
			Expect(updatedVm.Spec.Template.Spec.Volumes).To(Equal(vm.Spec.Template.Spec.Volumes), "volumes should be unchanged")
			Expect(updatedVm.Spec.DataVolumeTemplates).To(Equal(vm.Spec.DataVolumeTemplates), "datavolumeTemplates should be unchanged")

			memory := vm.Spec.Template.Spec.Domain.Resources.Requests[k8sv1.ResourceMemory]
			memoryUpdated := updatedVm.Spec.Template.Spec.Domain.Resources.Requests[k8sv1.ResourceMemory]
			Expect(memoryUpdated.Value()).To(Equal(memory.Value()), "memory should be unchanged")

		})

		It("taskrun succeded and disks are deleted and replaced by a new one", func() {
			config := &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:                testtemplate.NewRhelDesktopTinyTemplate().Build(),
					TemplateName:            testtemplate.RhelTemplateName,
					SourceTemplateNamespace: DeployTargetNS,
					DeleteDisks:             true,
					Disks:                   MockDisk,
				},
			}
			f.TestSetup(config)

			template := config.TaskData.Template
			if template != nil {
				template, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(template)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.TemplateName,
					"namespace": config.TaskData.TemplateNamespace,
				})

			updatedTemplate, err := f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.Background(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(updatedTemplate).ToNot(BeNil(), "new template should exists")
			f.ManageTemplates(updatedTemplate)

			updatedVm, _, err := zutils.DecodeVM(updatedTemplate)
			Expect(err).ShouldNot(HaveOccurred(), "decode updated VM")

			Expect(len(updatedVm.Spec.Template.Spec.Domain.Devices.Disks)).To(Equal(1), "disks should equal")
			Expect(updatedVm.Spec.Template.Spec.Domain.Devices.Disks[0].Name).To(Equal(Disk.Name), "disk name should equal")
		})

		It("taskrun succeded and volumes are deleted and replaced by a new one", func() {
			config := &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:                testtemplate.NewRhelDesktopTinyTemplate().Build(),
					TemplateName:            testtemplate.RhelTemplateName,
					SourceTemplateNamespace: DeployTargetNS,
					DeleteVolumes:           true,
					Volumes:                 MockVolume,
				},
			}
			f.TestSetup(config)

			template := config.TaskData.Template
			if template != nil {
				template, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(template)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.TemplateName,
					"namespace": config.TaskData.TemplateNamespace,
				})

			updatedTemplate, err := f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.Background(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(updatedTemplate).ToNot(BeNil(), "new template should exists")
			f.ManageTemplates(updatedTemplate)

			updatedVm, _, err := zutils.DecodeVM(updatedTemplate)
			Expect(err).ShouldNot(HaveOccurred(), "decode updated VM")

			Expect(len(updatedVm.Spec.Template.Spec.Volumes)).To(Equal(1), "Volumes should equal")
			Expect(updatedVm.Spec.Template.Spec.Volumes[0].Name).To(Equal(Volume.Name), "Volumes should equal")
		})

		It("taskrun succeded and template parameters are deleted and replaced by a new one", func() {
			config := &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:                 testtemplate.NewRhelDesktopTinyTemplate().Build(),
					TemplateName:             testtemplate.RhelTemplateName,
					SourceTemplateNamespace:  DeployTargetNS,
					DeleteTemplateParameters: true,
					TemplateParameters:       MockTemplateParameter,
				},
			}
			f.TestSetup(config)

			template := config.TaskData.Template
			if template != nil {
				template, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(template)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.TemplateName,
					"namespace": config.TaskData.TemplateNamespace,
				})

			updatedTemplate, err := f.TemplateClient.Templates(string(config.TaskData.TemplateNamespace)).Get(context.Background(), config.TaskData.TemplateName, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(updatedTemplate).ToNot(BeNil(), "new template should exists")
			f.ManageTemplates(updatedTemplate)

			Expect(len(updatedTemplate.Parameters)).To(Equal(1), "parameters should equal")
			Expect(updatedTemplate.Parameters[0].Name).To(Equal(TemplateParameters[0].Name), "parameter name should equal")
		})

		It("taskrun succeeded and template was deleted", func() {
			config := &testconfigs.ModifyTemplateTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: ModifyTemplateServiceAccountName,
				},
				TaskData: testconfigs.ModifyTemplateTaskData{
					Template:                testtemplate.NewRhelDesktopTinyTemplate().Build(),
					TemplateName:            testtemplate.RhelTemplateName,
					SourceTemplateNamespace: DeployTargetNS,
					DeleteTemplate:          true,
				},
			}
			f.TestSetup(config)

			template := config.TaskData.Template
			if template != nil {
				template, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(template)
			}

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					"name":      config.TaskData.TemplateName,
					"namespace": config.TaskData.TemplateNamespace,
				})

			Eventually(func() bool {
				_, err := f.TemplateClient.Templates(config.TaskData.TemplateNamespace).Get(context.Background(), config.TaskData.TemplateName, v1.GetOptions{})
				return errors.IsNotFound(err)
			}, time.Second*360, time.Second).Should(BeTrue())
		})
	})
})
