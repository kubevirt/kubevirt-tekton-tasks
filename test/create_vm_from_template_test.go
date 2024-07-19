package test

import (
	"context"

	testtemplate "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/testconfigs"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/vm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Create VM from template", func() {
	f := framework.NewFramework().
		LimitEnvScope(OKDEnvScope).
		OnBeforeTestSetup(func(config framework.TestConfig) {
			if createVMConfig, ok := config.(*testconfigs.CreateVMTestConfig); ok {
				createVMConfig.TaskData.CreateMode = CreateVMTemplateMode
			}
		})

	DescribeTable("taskrun fails and no VM is created", func(config *testconfigs.CreateVMTestConfig) {
		f.TestSetup(config)

		if template := config.TaskData.Template; template != nil {
			template, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(template)
		}

		expectedVM := config.TaskData.GetExpectedVMStubMeta()
		f.ManageVMs(expectedVM) // in case it succeeds

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectFailure().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResults(nil)

		_, err := vm.WaitForVM(f.KubevirtClient, expectedVM.Namespace, expectedVM.Name,
			"", config.GetTaskRunTimeout(), false)
		Expect(err).Should(HaveOccurred())
	},
		Entry("no template specified", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "only one of vm-manifest, template-name or virtctl should be specified",
			},
			TaskData: testconfigs.CreateVMTaskData{},
		}),
		Entry("template with no VM", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "no VM object found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().WithNoVM().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("invalid-template-no-vm")),
				},
			},
		}),
		Entry("non existent template", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "templates.template.openshift.io \"non-existent-template\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				TemplateName: "non-existent-template",
			},
		}),
		Entry("invalid template params", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "invalid template-params: no key found before \"InvalidDescription\"; pair should be in \"KEY:VAL\" format",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().WithDescriptionParam().Build(),
				TemplateParams: []string{
					"InvalidDescription",
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("invalid-template-params")),
				},
			},
		}),
		Entry("missing template params", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "required params are missing values: NAME",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
			},
		}),

		Entry("missing one template param", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "required params are missing values: DESCRIPTION",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().WithDescriptionParam().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("missing-one-template-param")),
				},
			},
		}),
		Entry("create vm with non matching disk fails", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "admission webhook \"virtualmachine-validator.kubevirt.io\" denied the request: spec.template.spec.domain.devices.disks[0].Name",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().WithNonMatchingDisk().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("non-matching-disk-fail-creation")),
				},
			},
		}),
	)

	DescribeTable("VM is created successfully", func(config *testconfigs.CreateVMTestConfig) {
		f.TestSetup(config)
		if template := config.TaskData.Template; template != nil {
			template, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(template)
		}

		expectedVM := config.TaskData.GetExpectedVMStubMeta()
		f.ManageVMs(expectedVM)

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResults(map[string]string{
				CreateVMResults.Name:      expectedVM.Name,
				CreateVMResults.Namespace: expectedVM.Namespace,
			})

		vm, err := vm.WaitForVM(f.KubevirtClient, expectedVM.Namespace, expectedVM.Name,
			"", config.GetTaskRunTimeout(), false)
		Expect(err).ShouldNot(HaveOccurred())

		if config.TaskData.SetOwnerReference == "true" {
			Expect(vm.OwnerReferences).To(HaveLen(1), "vm should has owner reference")
			Expect(vm.OwnerReferences[0].Kind).To(Equal("Pod"), "OwnerReference should have Kind Pod")
			Expect(vm.OwnerReferences[0].Name).To(HavePrefix("e2e-tests-taskrun-vm-create"), "OwnerReference should be binded to correct Pod")
		} else {
			Expect(vm.OwnerReferences).To(BeEmpty(), "vm OwnerReference should be empty")
		}
	},
		Entry("simple vm", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("simple-vm")),
				},
				SetOwnerReference: "true",
			},
		}),
		Entry("vm to deploy namespace by default", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-to-deploy-by-default")),
				},
				SetOwnerReference: "false",
			},
		}),
		Entry("vm with template from deploy namespace by default", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-template-from-deploy-by-default")),
				},
			},
		}),
		Entry("vm with template to deploy namespace by default", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-template-to-deploy-by-default")),
				},
			},
		}),
		Entry("vm with multiple params with template in deploy NS", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "description: e2e description with spaces",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().WithDescriptionParam().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.DescriptionParam, "e2e description with spaces"),
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-params")),
					testtemplate.TemplateParam("test", "test test"),
				},
			},
		}),
		Entry("works also in the same namespace as deploy", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("same-ns")),
				},
			},
		}),
	)

	It("VM from common template is created successfully + test trim", func() {
		config := &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				TemplateName:      SpacesSmall + "fedora-server-medium",
				TemplateNamespace: SpacesSmall + "openshift" + SpacesSmall,
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-common-template")),
				},
				IsCommonTemplate: true,
			},
		}
		f.TestSetup(config)

		expectedVM := config.TaskData.GetExpectedVMStubMeta()
		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResults(map[string]string{
				CreateVMResults.Name:      expectedVM.Name,
				CreateVMResults.Namespace: expectedVM.Namespace,
			})

		// don't wait for the DV (just VM creation), because it could take a long time and the size is pretty big
		_, err := vm.WaitForVM(f.KubevirtClient, expectedVM.Namespace, expectedVM.Name,
			"", config.GetTaskRunTimeout(), true)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("VM is created from template properly", func() {
		template := testtemplate.NewCirrosServerTinyTemplate().Build()
		config := &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: template,
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-template-data")),
				},
			},
		}
		f.TestSetup(config)
		template, err := f.TemplateClient.Templates(template.Namespace).Create(context.Background(), template, v1.CreateOptions{})
		Expect(err).ShouldNot(HaveOccurred())
		f.ManageTemplates(template)

		expectedVMStub := config.TaskData.GetExpectedVMStubMeta()
		f.ManageVMs(expectedVMStub)

		runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(config.GetAllExpectedLogs()...).
			ExpectResults(map[string]string{
				CreateVMResults.Name:      expectedVMStub.Name,
				CreateVMResults.Namespace: expectedVMStub.Namespace,
			})

		vm, err := vm.WaitForVM(f.KubevirtClient, expectedVMStub.Namespace, expectedVMStub.Name,
			"", config.GetTaskRunTimeout(), false)
		Expect(err).ShouldNot(HaveOccurred())

		vmName := expectedVMStub.Name
		expectedVM := testtemplate.GetVM(template)
		// fill template VM accordingly
		expectedVM.Spec.Template.Spec.Hostname = vmName
		expectedVM.Spec.Template.Spec.Domain.Machine = vm.Spec.Template.Spec.Domain.Machine // ignore Machine
		expectedVM.Spec.Template.Spec.Architecture = vm.Spec.Template.Spec.Architecture     // ignore Architecture
		expectedVM.Spec.Template.ObjectMeta.Labels["vm.kubevirt.io/name"] = vm.Spec.Template.ObjectMeta.Name

		Expect(vm.Spec.Template.Spec).Should(Equal(expectedVM.Spec.Template.Spec))
		// check VM labels
		Expect(vm.Labels).Should(Equal(map[string]string{
			"app":                              vmName,
			"vm.kubevirt.io/template":          "centos-server-tiny",
			"vm.kubevirt.io/template.revision": "147",
			"vm.kubevirt.io/template.version":  "0.3.2",
		}))
		// check VMI labels
		Expect(vm.Spec.Template.ObjectMeta.Labels).Should(Equal(map[string]string{
			"kubevirt.io/domain": vmName,
			"kubevirt.io/size":   "tiny",
		}))
	})

	Context("with StartVM", func() {
		DescribeTable("VM is created from template with StartVM attribute", func(config *testconfigs.CreateVMTestConfig, phase kubevirtv1.VirtualMachineInstancePhase, running bool) {
			f.TestSetup(config)
			template, err := f.TemplateClient.Templates(config.TaskData.Template.Namespace).Create(context.Background(), config.TaskData.Template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(template)

			expectedVMStub := config.TaskData.GetExpectedVMStubMeta()
			f.ManageVMs(expectedVMStub)

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					CreateVMResults.Name:      expectedVMStub.Name,
					CreateVMResults.Namespace: expectedVMStub.Namespace,
				})

			vm, err := vm.WaitForVM(f.KubevirtClient, expectedVMStub.Namespace, expectedVMStub.Name,
				phase, config.GetTaskRunTimeout(), false)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*vm.Spec.Running).To(Equal(running), "vm should be in correct running phase")
		},
			Entry("with invalid StartVM value", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateParams: []string{
						testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-template-data")),
					},
					StartVM: "invalid_value",
				},
			}, kubevirtv1.VirtualMachineInstancePhase(""), false),
			Entry("with false StartVM value", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateParams: []string{
						testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-template-data")),
					},
					StartVM: "false",
				},
			}, kubevirtv1.VirtualMachineInstancePhase(""), false),
			Entry("with true StartVM value", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateParams: []string{
						testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-template-data")),
					},
					StartVM: "true",
				},
			}, kubevirtv1.Running, true),
		)
	})

	Context("with RunStrategy", func() {
		DescribeTable("VM is created from template with runStrategy attribute", func(config *testconfigs.CreateVMTestConfig, expectedRunStrategy kubevirtv1.VirtualMachineRunStrategy) {
			f.TestSetup(config)

			template, err := f.TemplateClient.Templates(config.TaskData.Template.Namespace).Create(context.Background(), config.TaskData.Template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(template)

			expectedVMStub := config.TaskData.GetExpectedVMStubMeta()
			f.ManageVMs(expectedVMStub)

			runner.NewTaskRunRunner(f, config.GetTaskRun()).
				CreateTaskRun().
				ExpectSuccess().
				ExpectLogs(config.GetAllExpectedLogs()...).
				ExpectResults(map[string]string{
					CreateVMResults.Name:      expectedVMStub.Name,
					CreateVMResults.Namespace: expectedVMStub.Namespace,
				})

			vm, err := f.KubevirtClient.VirtualMachine(expectedVMStub.Namespace).Get(context.Background(), expectedVMStub.Name, v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(*vm.Spec.RunStrategy).To(Equal(expectedRunStrategy), "vm should have correct run strategy")
		},
			Entry("with RunStrategy always", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateParams: []string{
						testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-template-data")),
					},
					RunStrategy: "Always",
					StartVM:     "true",
				},
			}, kubevirtv1.RunStrategyAlways),
			Entry("with RunStrategy halted", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateParams: []string{
						testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-template-data")),
					},
					RunStrategy: "Halted",
				},
			}, kubevirtv1.RunStrategyHalted),
			Entry("with RunStrategy halted and startVM", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateParams: []string{
						testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-template-data")),
					},
					RunStrategy: "Halted",
					StartVM:     "true",
				},
			}, kubevirtv1.RunStrategyAlways),
			Entry("with RunStrategy Manual", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateParams: []string{
						testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-template-data")),
					},
					RunStrategy: "Manual",
					StartVM:     "true",
				},
			}, kubevirtv1.RunStrategyManual),
			Entry("with RunStrategy RerunOnFailure", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ExpectedLogs: ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
					TemplateParams: []string{
						testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-template-data")),
					},
					RunStrategy: "RerunOnFailure",
					StartVM:     "true",
				},
			}, kubevirtv1.RunStrategyRerunOnFailure),
		)
	})
})
