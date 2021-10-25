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
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

var _ = Describe("Create VM from template", func() {
	f := framework.NewFramework().
		LimitEnvScope(OKDEnvScope).
		OnBeforeTestSetup(func(config framework.TestConfig) {
			if createVMConfig, ok := config.(*testconfigs.CreateVMTestConfig); ok {
				createVMConfig.TaskData.CreateMode = CreateVMTemplateMode
			}
		})

	table.DescribeTable("taskrun fails and no VM is created", func(config *testconfigs.CreateVMTestConfig) {
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

		_, err := vm.WaitForVM(f.KubevirtClient, f.CdiClient, expectedVM.Namespace, expectedVM.Name,
			"", config.GetTaskRunTimeout(), false)
		Expect(err).Should(HaveOccurred())
	},
		table.Entry("no service account", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "cannot get resource \"templates\" in API group \"template.openshift.io\"",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("no-sa")),
				},
			},
		}),
		table.Entry("no template specified", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "one of vm-manifest, template-name should be specified",
			},
			TaskData: testconfigs.CreateVMTaskData{},
		}),
		table.Entry("template with no VM", &testconfigs.CreateVMTestConfig{
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
		table.Entry("non existent template", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "templates.template.openshift.io \"non-existent-template\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				TemplateTargetNamespace: TestTargetNS,
				TemplateName:            "non-existent-template",
			},
		}),
		table.Entry("invalid template params", &testconfigs.CreateVMTestConfig{
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
		table.Entry("missing template params", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "required params are missing values: NAME",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
			},
		}),

		table.Entry("missing one template param", &testconfigs.CreateVMTestConfig{
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
		table.Entry("non existent dv", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "datavolumes.cdi.kubevirt.io \"non-existent-dv\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-non-existent-dv")),
				},
				DataVolumes: []string{"non-existent-dv"},
			},
		}),
		table.Entry("non existent owned dv", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromTemplateServiceAccountName,
				ExpectedLogs:   "datavolumes.cdi.kubevirt.io \"non-existent-own-dv\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				Template: testtemplate.NewCirrosServerTinyTemplate().Build(),
				TemplateParams: []string{
					testtemplate.TemplateParam(testtemplate.NameParam, E2ETestsRandomName("vm-with-non-existent-owned-dv")),
				},
				OwnDataVolumes: []string{"non-existent-own-dv"},
			},
		}),
		table.Entry("non existent pvc", &testconfigs.CreateVMTestConfig{
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
		table.Entry("non existent owned pvcs", &testconfigs.CreateVMTestConfig{
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
		table.Entry("create vm with non matching disk fails", &testconfigs.CreateVMTestConfig{
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
		table.Entry("[NAMESPACE SCOPED] cannot create a VM in different namespace", &testconfigs.CreateVMTestConfig{
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
		table.Entry("[NAMESPACE SCOPED] cannot use template from a different namespace", &testconfigs.CreateVMTestConfig{
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

	table.DescribeTable("VM is created successfully", func(config *testconfigs.CreateVMTestConfig) {
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

		_, err := vm.WaitForVM(f.KubevirtClient, f.CdiClient, expectedVM.Namespace, expectedVM.Name,
			"", config.GetTaskRunTimeout(), false)
		Expect(err).ShouldNot(HaveOccurred())
	},
		table.Entry("simple vm", &testconfigs.CreateVMTestConfig{
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
		table.Entry("vm to deploy namespace by default", &testconfigs.CreateVMTestConfig{
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
		table.Entry("vm with template from deploy namespace by default", &testconfigs.CreateVMTestConfig{
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
		table.Entry("vm with template to deploy namespace by default", &testconfigs.CreateVMTestConfig{
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
		table.Entry("vm with multiple params with template in deploy NS", &testconfigs.CreateVMTestConfig{
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
		table.Entry("[CLUSTER SCOPED] works also in the same namespace as deploy", &testconfigs.CreateVMTestConfig{
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
				IsCommonTemplate: true,
				DataVolumesToCreate: []*datavolume.TestDataVolume{
					datavolume.NewBlankDataVolume("common-templates-src-dv"),
				},
				StartVM: "false",
			},
		}
		f.TestSetup(config)

		expectedVM := config.TaskData.GetExpectedVMStubMeta()
		f.ManageVMs(expectedVM)

		for _, dvWrapper := range config.TaskData.DataVolumesToCreate {
			dataVolume, err := f.CdiClient.DataVolumes(dvWrapper.Data.Namespace).Create(context.TODO(), dvWrapper.Data, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataVolumes(dataVolume)
			config.TaskData.TemplateParams = append(config.TaskData.TemplateParams, testtemplate.TemplateParam(SpacesSmall+testtemplate.SrcPvcNameParam, dataVolume.Name))
			config.TaskData.TemplateParams = append(config.TaskData.TemplateParams, testtemplate.TemplateParam(SpacesSmall+testtemplate.SrcPvcNamespaceParam, dataVolume.Namespace))
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
		_, err := vm.WaitForVM(f.KubevirtClient, f.CdiClient, expectedVM.Namespace, expectedVM.Name,
			"", config.GetTaskRunTimeout(), true)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("VM is created from template properly and running", func() {
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
				StartVM: "true",
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

		vm, err := vm.WaitForVM(f.KubevirtClient, f.CdiClient, expectedVMStub.Namespace, expectedVMStub.Name,
			kubevirtv1.Running, config.GetTaskRunTimeout(), false)
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
})
