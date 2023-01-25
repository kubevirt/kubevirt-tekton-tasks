package test

import (
	"context"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datavolume"
	testtemplate "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/vm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	cdiv1beta12 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
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
			template, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
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
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "one of vm-manifest, template-name should be specified",
			},
			TaskData: testconfigs.CreateVMTaskData{},
		}),
		Entry("template with no VM", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "no VM object found",
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
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "templates.template.openshift.io \"non-existent-template\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				TemplateTargetNamespace: TestTargetNS,
				TemplateName:            "non-existent-template",
			},
		}),
		Entry("invalid template params", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "invalid template-params: no key found before \"InvalidDescription\"; pair should be in \"KEY:VAL\" format",
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
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "required params are missing values: NAME",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
			},
		}),

		Entry("missing one template param", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "required params are missing values: DESCRIPTION",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().WithDescriptionParam().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("missing-one-template-param")),
				},
			},
		}),
		Entry("non existent dv", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "persistentvolumeclaims \"non-existent-dv\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-non-existent-dv")),
				},
				DataVolumes: []string{"non-existent-dv"},
			},
		}),
		Entry("non existent owned dv", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "persistentvolumeclaims \"non-existent-own-dv\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-non-existent-owned-dv")),
				},
				OwnDataVolumes: []string{"non-existent-own-dv"},
			},
		}),
		Entry("non existent pvc", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "persistentvolumeclaims \"non-existent-pvc\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-non-existent-pvc")),
				},
				PersistentVolumeClaims: []string{"non-existent-pvc"},
			},
		}),
		Entry("non existent owned pvcs", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "persistentvolumeclaims \"non-existent-own-pvc\" not found\npersistentvolumeclaims \"non-existent-own-pvc-2\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-non-existent-owned-pvcs")),
				},
				OwnPersistentVolumeClaims: []string{"non-existent-own-pvc", "non-existent-own-pvc-2"},
			},
		}),
		Entry("create vm with non matching disk fails", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "admission webhook \"virtualmachine-validator.kubevirt.io\" denied the request: spec.template.spec.domain.devices.disks[0].Name",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().WithNonMatchingDisk().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("non-matching-disk-fail-creation")),
				},
			},
		}),
		Entry("[NAMESPACE SCOPED] cannot create a VM in different namespace", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "processedtemplates.template.openshift.io is forbidden",
				LimitTestScope: NamespaceTestScope,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("different-ns-namespace-scope")),
				},
				VMTargetNamespace: SystemTargetNS,
			},
		}),
		Entry("[NAMESPACE SCOPED] cannot use template from a different namespace", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "templates.template.openshift.io \"unreachable-template\" is forbidden",
				LimitTestScope: NamespaceTestScope,
			},
			TaskData: testconfigs.CreateVMTaskData{
				TemplateTargetNamespace: SystemTargetNS,
				TemplateName:            "unreachable-template",
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("different-ns-namespace-scope")),
				},
			},
		}),
	)

	DescribeTable("VM is created successfully", func(config *testconfigs.CreateVMTestConfig) {
		f.TestSetup(config)
		if template := config.TaskData.Template; template != nil {
			template, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
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

		_, err := vm.WaitForVM(f.KubevirtClient, expectedVM.Namespace, expectedVM.Name,
			"", config.GetTaskRunTimeout(), false)
		Expect(err).ShouldNot(HaveOccurred())
	},
		Entry("simple vm", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("simple-vm")),
				},
			},
		}),
		Entry("vm to deploy namespace by default", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-to-deploy-by-default")),
				},
				VMTargetNamespace:                  DeployTargetNS,
				UseDefaultVMNamespacesInTaskParams: true,
			},
		}),
		Entry("vm with template from deploy namespace by default", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-template-from-deploy-by-default")),
				},
				TemplateTargetNamespace:                  DeployTargetNS,
				UseDefaultTemplateNamespacesInTaskParams: true,
			},
		}),
		Entry("vm with template to deploy namespace by default", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-template-to-deploy-by-default")),
				},
				VMTargetNamespace:                        DeployTargetNS,
				TemplateTargetNamespace:                  DeployTargetNS,
				UseDefaultVMNamespacesInTaskParams:       true,
				UseDefaultTemplateNamespacesInTaskParams: true,
			},
		}),
		Entry("vm with multiple params with template in deploy NS", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "description: e2e description with spaces",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template:                testtemplate.NewCirrosServerTinyTemplate().WithDescriptionParam().Build(),
				TemplateTargetNamespace: DeployTargetNS,
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.DescriptionParam, "e2e description with spaces"),
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-params")),
					testtemplate.TemplateParam("test", "test test"),
				},
			},
		}),
		Entry("[CLUSTER SCOPED] works also in the same namespace as deploy", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				LimitTestScope: ClusterTestScope,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template:                testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateTargetNamespace: DeployTargetNS,
				VMTargetNamespace:       DeployTargetNS,
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("same-ns-cluster-scope")),
				},
			},
		}),
	)

	It("VM from common template is created successfully + test trim", func() {
		config := &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				TemplateName:      SpacesSmall + "fedora-server-tiny",
				TemplateNamespace: SpacesSmall + "openshift" + SpacesSmall,
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-common-template")),
				},
				IsCommonTemplate:  true,
				VMTargetNamespace: DeployTargetNS,
				DataVolumesToCreate: []*datavolume.TestDataVolume{
					datavolume.NewBlankDataVolume("common-templates-src-dv"),
				},
			},
		}
		f.TestSetup(config)

		expectedVM := config.TaskData.GetExpectedVMStubMeta()
		f.ManageVMs(expectedVM)

		for _, dvWrapper := range config.TaskData.DataVolumesToCreate {
			dataVolume, err := f.CdiClient.DataVolumes(dvWrapper.Data.Namespace).Create(context.TODO(), dvWrapper.Data, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataVolumes(dataVolume)

			datasource := cdiv1beta12.DataSource{
				ObjectMeta: v1.ObjectMeta{
					Name:      dvWrapper.Data.Name,
					Namespace: dvWrapper.Data.Namespace,
				},
				Spec: cdiv1beta12.DataSourceSpec{
					Source: cdiv1beta12.DataSourceSource{
						PVC: &cdiv1beta12.DataVolumeSourcePVC{
							Name:      dvWrapper.Data.Name,
							Namespace: dvWrapper.Data.Namespace,
						},
					},
				},
			}
			ds, err := f.CdiClient.DataSources(dvWrapper.Data.Namespace).Create(context.TODO(), &datasource, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataSources(ds)

			config.TaskData.TemplateParams = append(config.TaskData.TemplateParams, testtemplate.TemplateParam(SpacesSmall+testtemplate.DataVolumeNameParam, dataVolume.Name))
			config.TaskData.TemplateParams = append(config.TaskData.TemplateParams, testtemplate.TemplateParam(SpacesSmall+testtemplate.DataVolumeNamespaceParam, dataVolume.Namespace))
		}

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
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: template,
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-from-template-data")),
				},
			},
		}
		f.TestSetup(config)
		template, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
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

	Context("with StartVM", func() {
		DescribeTable("VM is created from template with StartVM attribute", func(config *testconfigs.CreateVMTestConfig, phase kubevirtv1.VirtualMachineInstancePhase, running bool) {
			f.TestSetup(config)
			template, err := f.TemplateClient.Templates(config.TaskData.Template.Namespace).Create(context.TODO(), config.TaskData.Template, v1.CreateOptions{})
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
					ServiceAccount: CreateVMFromTemplateServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
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
					ServiceAccount: CreateVMFromTemplateServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
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
					ServiceAccount: CreateVMFromTemplateServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
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

			template, err := f.TemplateClient.Templates(config.TaskData.Template.Namespace).Create(context.TODO(), config.TaskData.Template, v1.CreateOptions{})
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

			vm, err := f.KubevirtClient.VirtualMachine(expectedVMStub.Namespace).Get(expectedVMStub.Name, &v1.GetOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(*vm.Spec.RunStrategy).To(Equal(expectedRunStrategy), "vm should have correct run strategy")
		},
			Entry("with RunStrategy always", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CreateVMFromTemplateServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
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
					ServiceAccount: CreateVMFromTemplateServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
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
					ServiceAccount: CreateVMFromTemplateServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
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
					ServiceAccount: CreateVMFromTemplateServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
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
					ServiceAccount: CreateVMFromTemplateServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
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
