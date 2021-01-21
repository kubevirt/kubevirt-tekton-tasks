package test

import (
	"fmt"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/dv"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	templ "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/template"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/vm"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

const spaces = "  "

var _ = Describe("Create VM from template", func() {
	f := framework.NewFramework().LimitEnvScope(OpenshiftEnvScope)

	table.DescribeTable("taskrun fails and no VM is created", func(config *testconfigs.CreateVMFromTemplateTestConfig) {
		f.TestSetup(config)

		if template := config.TaskData.Template; template != nil {
			template, err := f.TemplateClient.Templates(template.Namespace).Create(template)
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(template)
		}

		expectedVM := config.TaskData.GetExpectedVMStubMeta()
		f.ManageVMs(expectedVM) // in case it succeeds

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectFailure().
			ExpectLogs(config.ExpectedLogs).
			ExpectResults(nil)

		_, err := vm.WaitForVM(f.KubevirtClient, f.CdiClient, expectedVM.Namespace, expectedVM.Name,
			"", config.GetTaskRunTimeout(), false)
		Expect(err).Should(HaveOccurred())
	},
		table.Entry("no service account", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "cannot get resource \"templates\" in API group \"template.openshift.io\"",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("no-sc")),
				},
			},
		}),
		table.Entry("no template specified", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "resource name may not be empty",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{},
		}),
		table.Entry("template with no VM", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "no VM object found",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().WithNoVM().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("invalid-template-no-vm")),
				},
			},
		}),
		table.Entry("non existent template", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "templates.template.openshift.io \"non-existent-template\" not found",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				TemplateTargetNamespace: TestTargetNS,
				TemplateName:            "non-existent-template",
			},
		}),
		table.Entry("invalid template params", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "parameters have incorrect format",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().WithDescriptionParam().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("invalid-template-params")),
					"InvalidDescription",
				},
			},
		}),
		table.Entry("missing template params", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "required params are missing values: NAME",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
			},
		}),

		table.Entry("missing one template param", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "required params are missing values: DESCRIPTION",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().WithDescriptionParam().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("missing-one-template-param")),
				},
			},
		}),
		table.Entry("non existent dv", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "datavolumes.cdi.kubevirt.io \"non-existent-dv\" not found",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("vm-with-non-existent-dv")),
				},
				DataVolumes: []string{"non-existent-dv"},
			},
		}),
		table.Entry("non existent owned dv", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "datavolumes.cdi.kubevirt.io \"non-existent-own-dv\" not found",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("vm-with-non-existent-owned-dv")),
				},
				OwnDataVolumes: []string{"non-existent-own-dv"},
			},
		}),
		table.Entry("non existent pvc", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "persistentvolumeclaims \"non-existent-pvc\" not found",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("vm-with-non-existent-pvc")),
				},
				PersistentVolumeClaims: []string{"non-existent-pvc"},
			},
		}),
		table.Entry("non existent owned pvcs", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "persistentvolumeclaims \"non-existent-own-pvc\" not found\npersistentvolumeclaims \"non-existent-own-pvc-2\" not found",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("vm-with-non-existent-owned-pvcs")),
				},
				OwnPersistentVolumeClaims: []string{"non-existent-own-pvc", "non-existent-own-pvc-2"},
			},
		}),
		table.Entry("create vm with non matching disk fails", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "admission webhook \"virtualmachine-validator.kubevirt.io\" denied the request: spec.template.spec.domain.devices.disks[0].Name",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().WitNonMatchingDisk().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("non-matching-disk-fail-creation")),
				},
			},
		}),
		table.Entry("[NAMESPACE SCOPED] cannot create a VM in different namespace", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "processedtemplates.template.openshift.io is forbidden",
				LimitTestScope: NamespaceTestScope,
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("different-ns-namespace-scope")),
				},
				VMTargetNamespace: SystemTargetNS,
			},
		}),
		table.Entry("[NAMESPACE SCOPED] cannot use template from a different namespace", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "templates.template.openshift.io \"unreachable-template\" is forbidden",
				LimitTestScope: NamespaceTestScope,
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				TemplateTargetNamespace: SystemTargetNS,
				TemplateName:            "unreachable-template",
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("different-ns-namespace-scope")),
				},
			},
		}),
	)

	table.DescribeTable("VM is created successfully", func(config *testconfigs.CreateVMFromTemplateTestConfig) {
		f.TestSetup(config)
		if template := config.TaskData.Template; template != nil {
			template, err := f.TemplateClient.Templates(template.Namespace).Create(template)
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(template)
		}

		expectedVM := config.TaskData.GetExpectedVMStubMeta()
		f.ManageVMs(expectedVM)

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(config.ExpectedLogs).
			ExpectResults(map[string]string{
				CreateVMFromManifestResults.Name:      expectedVM.Name,
				CreateVMFromManifestResults.Namespace: expectedVM.Namespace,
			})

		_, err := vm.WaitForVM(f.KubevirtClient, f.CdiClient, expectedVM.Namespace, expectedVM.Name,
			"", config.GetTaskRunTimeout(), false)
		Expect(err).ShouldNot(HaveOccurred())
	},
		table.Entry("simple vm", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("simple-vm")),
				},
			},
		}),
		table.Entry("vm to deploy namespace by default", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("vm-to-deploy-by-default")),
				},
				VMTargetNamespace:                  DeployTargetNS,
				UseDefaultVMNamespacesInTaskParams: true,
			},
		}),
		table.Entry("vm with template from deploy namespace by default", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("vm-with-template-from-deploy-by-default")),
				},
				TemplateTargetNamespace:                  DeployTargetNS,
				UseDefaultTemplateNamespacesInTaskParams: true,
			},
		}),
		table.Entry("vm with template to deploy namespace by default", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: templ.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("vm-with-template-to-deploy-by-default")),
				},
				VMTargetNamespace:                        DeployTargetNS,
				TemplateTargetNamespace:                  DeployTargetNS,
				UseDefaultVMNamespacesInTaskParams:       true,
				UseDefaultTemplateNamespacesInTaskParams: true,
			},
		}),
		table.Entry("vm with multiple params with template in deploy NS", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "description: e2e-description",
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template:                templ.NewCirrosServerTinyTemplate().WithDescriptionParam().Build(),
				TemplateTargetNamespace: DeployTargetNS,
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("vm-with-params")),
					templ.TemplateParam(templ.DescriptionParam, "e2e-description"),
				},
			},
		}),
		table.Entry("[CLUSTER SCOPED] works also in the same namespace as deploy", &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				LimitTestScope: ClusterTestScope,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template:                templ.NewCirrosServerTinyTemplate().Build(),
				TemplateTargetNamespace: DeployTargetNS,
				VMTargetNamespace:       DeployTargetNS,
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("same-ns-cluster-scope")),
				},
			},
		}),
	)

	It("VM from common template is created successfully + test trim", func() {
		config := &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				TemplateName:      spaces + "fedora-server-tiny",
				TemplateNamespace: spaces + "openshift" + spaces,
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("vm-from-common-template")),
				},
				IsCommonTemplate: true,
				DataVolumesToCreate: []*dv.TestDataVolume{
					dv.NewBlankDataVolume("common-templates-src-dv"),
				},
			},
		}
		f.TestSetup(config)

		expectedVM := config.TaskData.GetExpectedVMStubMeta()
		f.ManageVMs(expectedVM)

		for _, dvWrapper := range config.TaskData.DataVolumesToCreate {
			dataVolume, err := f.CdiClient.DataVolumes(dvWrapper.Data.Namespace).Create(dvWrapper.Data)
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataVolumes(dataVolume)
			config.TaskData.TemplateParams = append(config.TaskData.TemplateParams, templ.TemplateParam(spaces+templ.SrcPvcNameParam, dataVolume.Name))
			config.TaskData.TemplateParams = append(config.TaskData.TemplateParams, templ.TemplateParam(spaces+templ.SrcPvcNamespaceParam, dataVolume.Namespace))
		}

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(config.ExpectedLogs).
			ExpectResults(map[string]string{
				CreateVMFromManifestResults.Name:      expectedVM.Name,
				CreateVMFromManifestResults.Namespace: expectedVM.Namespace,
			})

		// don't wait for the DV (just VM creation), because it could take a long time and the size is pretty big
		_, err := vm.WaitForVM(f.KubevirtClient, f.CdiClient, expectedVM.Namespace, expectedVM.Name,
			"", config.GetTaskRunTimeout(), true)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("VM with attached PVCs/DV is created successfully ", func() {
		runConfigurations := []map[dv.TestDataVolumeAttachmentType]int{
			{
				// try all at once
				dv.OwnedDV:  2,
				dv.OwnedPVC: 1,
				dv.PVC:      1,
				dv.DV:       1,
			},
		}

		// try for each type 1 or 2 dvs
		for count := 1; count < 3; count++ {
			for _, attachmentType := range []dv.TestDataVolumeAttachmentType{dv.OwnedDV, dv.OwnedPVC, dv.PVC, dv.DV} {
				runConfigurations = append(runConfigurations, map[dv.TestDataVolumeAttachmentType]int{
					attachmentType: count,
				})
			}
		}

		for idx, runConf := range runConfigurations {
			name := ""
			for attachmentType, count := range runConf {
				name += fmt.Sprintf("%v=%v ", attachmentType, count)
			}
			It(name, func() {
				var datavolumes []*dv.TestDataVolume
				for attachmentType, count := range runConf {
					name += fmt.Sprintf("%v=%v ", attachmentType, count)
					for id := 0; id < count; id++ {
						datavolumes = append(datavolumes,
							dv.NewBlankDataVolume(fmt.Sprintf("attach-to-vm-%v-%v", attachmentType, id)).AttachAs(attachmentType),
						)
					}
				}
				var expectedDisbBus string
				testTemplate := templ.NewCirrosServerTinyTemplate()
				switch idx % 4 { // try different disk buses for each test
				case 0:
					testTemplate.WithSataDiskValidations()
					expectedDisbBus = "sata"
				case 1:
					testTemplate.WithSCSIDiskValidations()
					expectedDisbBus = "scsi"
				case 2:
					testTemplate.WithVirtioDiskValidations()
					expectedDisbBus = "virtio"
				default:
					expectedDisbBus = "virtio"
				}
				config := &testconfigs.CreateVMFromTemplateTestConfig{
					TaskRunTestConfig: testconfigs.TaskRunTestConfig{
						ServiceAccount: CreateVMFromTemplateServiceAccountName,
						ExpectedLogs:   ExpectedSuccessfulVMCreation,
						Timeout:        Timeouts.SmallBlankDVCreation,
					},
					TaskData: testconfigs.CreateVMFromTemplateTaskData{
						Template: testTemplate.Build(),
						TemplateParams: []string{
							templ.TemplateParam(templ.NameParam, E2ETestsRandomName("simple-vm")),
						},
						DataVolumesToCreate:       datavolumes,
						ExpectedAdditionalDiskBus: expectedDisbBus,
					},
				}
				f.TestSetup(config)
				if template := config.TaskData.Template; template != nil {
					template, err := f.TemplateClient.Templates(template.Namespace).Create(template)
					Expect(err).ShouldNot(HaveOccurred())
					f.ManageTemplates(template)
				}
				for _, dvWrapper := range config.TaskData.DataVolumesToCreate {
					dataVolume, err := f.CdiClient.DataVolumes(dvWrapper.Data.Namespace).Create(dvWrapper.Data)
					Expect(err).ShouldNot(HaveOccurred())
					f.ManageDataVolumes(dataVolume)
					config.TaskData.SetDVorPVC(dataVolume.Name, dvWrapper.AttachmentType)
				}

				expectedVM := config.TaskData.GetExpectedVMStubMeta()
				f.ManageVMs(expectedVM)

				runner.NewTaskRunRunner(f, config.GetTaskRun()).
					CreateTaskRun().
					ExpectSuccess().
					ExpectLogs(config.ExpectedLogs).
					ExpectResults(map[string]string{
						CreateVMFromManifestResults.Name:      expectedVM.Name,
						CreateVMFromManifestResults.Namespace: expectedVM.Namespace,
					})

				vm, err := vm.WaitForVM(f.KubevirtClient, f.CdiClient, expectedVM.Namespace, expectedVM.Name,
					"", config.GetTaskRunTimeout(), false)
				Expect(err).ShouldNot(HaveOccurred())
				// check all disks are present
				Expect(vm.Spec.Template.Spec.Volumes).To(ConsistOf(expectedVM.Spec.Template.Spec.Volumes))
				Expect(vm.Spec.Template.Spec.Domain.Devices.Disks).To(ConsistOf(expectedVM.Spec.Template.Spec.Domain.Devices.Disks))
			})
		}
	})

	It("VM is created from template properly ", func() {
		template := templ.NewCirrosServerTinyTemplate().Build()
		config := &testconfigs.CreateVMFromTemplateTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMFromTemplateTaskData{
				Template: template,
				TemplateParams: []string{
					templ.TemplateParam(templ.NameParam, E2ETestsRandomName("vm-data")),
				},
			},
		}
		f.TestSetup(config)
		template, err := f.TemplateClient.Templates(template.Namespace).Create(template)
		Expect(err).ShouldNot(HaveOccurred())
		f.ManageTemplates(template)

		expectedVMStub := config.TaskData.GetExpectedVMStubMeta()
		f.ManageVMs(expectedVMStub)

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(config.ExpectedLogs).
			ExpectResults(map[string]string{
				CreateVMFromManifestResults.Name:      expectedVMStub.Name,
				CreateVMFromManifestResults.Namespace: expectedVMStub.Namespace,
			})

		vm, err := vm.WaitForVM(f.KubevirtClient, f.CdiClient, expectedVMStub.Namespace, expectedVMStub.Name,
			"", config.GetTaskRunTimeout(), false)
		Expect(err).ShouldNot(HaveOccurred())

		vmName := expectedVMStub.Name
		expectedVM := templ.GetVM(template)
		// fill template VM accordingly
		expectedVM.Spec.Template.Spec.Hostname = vmName
		expectedVM.Spec.Template.Spec.Domain.Machine = vm.Spec.Template.Spec.Domain.Machine // ignore Machine

		Expect(vm.Spec.Template.Spec).Should(Equal(expectedVM.Spec.Template.Spec))
		// check VM labels
		Expect(vm.Labels).Should(Equal(map[string]string{
			"app":                                  vmName,
			"flavor.template.kubevirt.io/tiny":     "true",
			"os.template.kubevirt.io/centos7.0":    "true",
			"workload.template.kubevirt.io/server": "true",
			"vm.kubevirt.io/template":              template.Name,
			"vm.kubevirt.io/template.namespace":    template.Namespace,
			"vm.kubevirt.io/template.revision":     "147",
			"vm.kubevirt.io/template.version":      "0.3.2",
		}))
		// check os annotation
		Expect(vm.Annotations["name.os.template.kubevirt.io/centos7.0"]).Should(Equal("CentOS 7.0 or higher"))
		// check VMI labels
		Expect(vm.Spec.Template.ObjectMeta.Labels).Should(Equal(map[string]string{
			"flavor.template.kubevirt.io/tiny":     "true",
			"os.template.kubevirt.io/centos7.0":    "true",
			"workload.template.kubevirt.io/server": "true",
			"kubevirt.io/domain":                   vmName,
			"kubevirt.io/size":                     "tiny",
			"vm.kubevirt.io/name":                  vmName,
		}))
	})
})
