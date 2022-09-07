package test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/vm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Create VM from manifest", func() {
	f := framework.NewFramework().
		OnBeforeTestSetup(func(config framework.TestConfig) {
			if createVMConfig, ok := config.(*testconfigs.CreateVMTestConfig); ok {
				createVMConfig.TaskData.CreateMode = CreateVMVMManifestMode
			}
		})

	BeforeEach(func() {
		if f.TestOptions.SkipCreateVMFromManifestTests {
			Skip("skipCreateVMFromManifestTests is set to true, skipping tests")
		}
	})

	DescribeTable("taskrun fails and no VM is created", func(config *testconfigs.CreateVMTestConfig) {
		f.TestSetup(config)

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
		Entry("no vm manifest", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ExpectedLogs: "one of vm-manifest, template-name should be specified",
			},
			TaskData: testconfigs.CreateVMTaskData{},
		}),
		Entry("invalid manifest", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   "could not read VM manifest: error unmarshaling",
			},
			TaskData: testconfigs.CreateVMTaskData{
				VMManifest: "invalid manifest",
			},
		}),
		Entry("non existent dv", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   "datavolumes.cdi.kubevirt.io \"non-existent-dv\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM:          testobjects.NewTestAlpineVM("vm-with-non-existent-dv").Build(),
				DataVolumes: []string{"non-existent-dv"},
			},
		}),
		Entry("non existent owned dv", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   "datavolumes.cdi.kubevirt.io \"non-existent-own-dv\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM:             testobjects.NewTestAlpineVM("vm-with-non-existent-owned-dv").Build(),
				OwnDataVolumes: []string{"non-existent-own-dv"},
			},
		}),
		Entry("non existent pvc", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   "persistentvolumeclaims \"non-existent-pvc\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM:                     testobjects.NewTestAlpineVM("vm-with-non-existent-pvc").Build(),
				PersistentVolumeClaims: []string{"non-existent-pvc"},
			},
		}),
		Entry("non existent owned pvcs", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   "persistentvolumeclaims \"non-existent-own-pvc\" not found\npersistentvolumeclaims \"non-existent-own-pvc-2\" not found",
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM:                        testobjects.NewTestAlpineVM("vm-with-non-existent-owned-pvcs").Build(),
				OwnPersistentVolumeClaims: []string{"non-existent-own-pvc", "non-existent-own-pvc-2"},
			},
		}),
		Entry("create vm with non matching disk fails", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   "admission webhook \"virtualmachine-validator.kubevirt.io\" denied the request: spec.template.spec.domain.devices.disks[0].Name",
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM: testobjects.NewTestAlpineVM("vm-with-non-existent-pvc").WithNonMatchingDisk().Build(),
			},
		}),
		Entry("[NAMESPACE SCOPED] cannot create a VM in different namespace", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   "cannot create resource \"virtualmachines\" in API group \"kubevirt.io\"",
				LimitTestScope: NamespaceTestScope,
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM:                testobjects.NewTestAlpineVM("different-ns-namespace-scope").Build(),
				VMTargetNamespace: SystemTargetNS,
			},
		}),
		Entry("[NAMESPACE SCOPED] cannot create a VM in different namespace in manifest", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   "cannot create resource \"virtualmachines\" in API group \"kubevirt.io\"",
				LimitTestScope: NamespaceTestScope,
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM:                        testobjects.NewTestAlpineVM("different-ns-namespace-scope-in-manifest").Build(),
				VMManifestTargetNamespace: SystemTargetNS,
			},
		}),
	)

	DescribeTable("VM is created successfully", func(config *testconfigs.CreateVMTestConfig) {
		f.TestSetup(config)

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
		Entry("simple vm", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM: testobjects.NewTestAlpineVM("simple-vm").Build(),
			},
		}),
		Entry("vm to deploy namespace by default", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM:                                 testobjects.NewTestAlpineVM("vm-to-deploy-by-default").Build(),
				VMTargetNamespace:                  DeployTargetNS,
				UseDefaultVMNamespacesInTaskParams: true,
			},
		}),
		Entry("vm with manifest namespace", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM:                                 testobjects.NewTestAlpineVM("vm-with-manifest-namespace").Build(),
				VMManifestTargetNamespace:          DeployTargetNS,
				UseDefaultVMNamespacesInTaskParams: true,
			},
		}),

		Entry("vm with overridden manifest namespace", &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM:                        testobjects.NewTestAlpineVM("vm-with-overridden-manifest-namespace").Build(),
				VMManifestTargetNamespace: DeployTargetNS,
			},
		}),
	)

	It("VM is created from manifest properly ", func() {
		config := &testconfigs.CreateVMTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{
				ServiceAccount: CreateVMFromManifestServiceAccountName,
				ExpectedLogs:   ExpectedSuccessfulVMCreation,
			},
			TaskData: testconfigs.CreateVMTaskData{
				VM: testobjects.NewTestAlpineVM("vm-from-manifest-data").
					WithLabel("app", "my-custom-app").
					WithVMILabel("name", "test").
					WithVMILabel("ra", "rara").
					Build(),
			},
		}
		f.TestSetup(config)

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
			"", config.GetTaskRunTimeout(), false)
		Expect(err).ShouldNot(HaveOccurred())

		vmName := expectedVMStub.Name
		expectedVM := config.TaskData.VM
		// fill VM accordingly
		expectedVM.Spec.Template.Spec.Domain.Machine = vm.Spec.Template.Spec.Domain.Machine // ignore Machine

		Expect(vm.Spec.Template.Spec).Should(Equal(expectedVM.Spec.Template.Spec))
		// check VM labels
		Expect(vm.Labels).Should(Equal(expectedVM.Labels))
		// check VMI labels
		Expect(vm.Spec.Template.ObjectMeta.Labels).Should(Equal(map[string]string{
			"name":                "test",
			"ra":                  "rara",
			"vm.kubevirt.io/name": vmName,
		}))
	})

	Context("with StartVM", func() {
		DescribeTable("VM is created successfully", func(config *testconfigs.CreateVMTestConfig, phase kubevirtv1.VirtualMachineInstancePhase, running bool) {
			f.TestSetup(config)

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
				phase, config.GetTaskRunTimeout(), false)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(*vm.Spec.Running).To(Equal(running), "vm should be in correct running phase")
		},
			Entry("with false StartVM value", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CreateVMFromManifestServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					VM: testobjects.NewTestAlpineVM("vm-from-manifest-data").
						WithLabel("app", "my-custom-app").
						WithVMILabel("name", "test").
						WithVMILabel("ra", "rara").
						Build(),
					StartVM: "false",
				},
			}, kubevirtv1.VirtualMachineInstancePhase(""), false),
			Entry("with invalid StartVM value", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CreateVMFromManifestServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					VM: testobjects.NewTestAlpineVM("vm-from-manifest-data").
						WithLabel("app", "my-custom-app").
						WithVMILabel("name", "test").
						WithVMILabel("ra", "rara").
						Build(),
					StartVM: "invalid_value",
				},
			}, kubevirtv1.VirtualMachineInstancePhase(""), false),
			Entry("with true StartVM value", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CreateVMFromManifestServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					VM: testobjects.NewTestAlpineVM("vm-from-manifest-data").
						WithLabel("app", "my-custom-app").
						WithVMILabel("name", "test").
						WithVMILabel("ra", "rara").
						Build(),
					StartVM: "true",
				},
			}, kubevirtv1.Running, true),
		)
	})

	Context("with RunStrategy", func() {
		DescribeTable("VM is created successfully", func(config *testconfigs.CreateVMTestConfig, expectedRunStrategy kubevirtv1.VirtualMachineRunStrategy) {
			f.TestSetup(config)

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
					ServiceAccount: CreateVMFromManifestServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					VM: testobjects.NewTestAlpineVM("vm-from-manifest-data").
						WithLabel("app", "my-custom-app").
						WithVMILabel("name", "test").
						WithVMILabel("ra", "rara").
						Build(),
					RunStrategy: "Always",
					StartVM:     "true",
				},
			}, kubevirtv1.RunStrategyAlways),
			Entry("with RunStrategy halted", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CreateVMFromManifestServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					VM: testobjects.NewTestAlpineVM("vm-from-manifest-data").
						WithLabel("app", "my-custom-app").
						WithVMILabel("name", "test").
						WithVMILabel("ra", "rara").
						Build(),
					RunStrategy: "Halted",
				},
			}, kubevirtv1.RunStrategyHalted),
			Entry("with RunStrategy halted and startVM", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CreateVMFromManifestServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					VM: testobjects.NewTestAlpineVM("vm-from-manifest-data").
						WithLabel("app", "my-custom-app").
						WithVMILabel("name", "test").
						WithVMILabel("ra", "rara").
						Build(),
					RunStrategy: "Halted",
					StartVM:     "true",
				},
			}, kubevirtv1.RunStrategyAlways),
			Entry("with RunStrategy Manual", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CreateVMFromManifestServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					VM: testobjects.NewTestAlpineVM("vm-from-manifest-data").
						WithLabel("app", "my-custom-app").
						WithVMILabel("name", "test").
						WithVMILabel("ra", "rara").
						Build(),
					RunStrategy: "Manual",
					StartVM:     "true",
				},
			}, kubevirtv1.RunStrategyManual),
			Entry("with RunStrategy RerunOnFailure", &testconfigs.CreateVMTestConfig{
				TaskRunTestConfig: testconfigs.TaskRunTestConfig{
					ServiceAccount: CreateVMFromManifestServiceAccountName,
					ExpectedLogs:   ExpectedSuccessfulVMCreation,
				},
				TaskData: testconfigs.CreateVMTaskData{
					VM: testobjects.NewTestAlpineVM("vm-from-manifest-data").
						WithLabel("app", "my-custom-app").
						WithVMILabel("name", "test").
						WithVMILabel("ra", "rara").
						Build(),
					RunStrategy: "RerunOnFailure",
					StartVM:     "true",
				},
			}, kubevirtv1.RunStrategyRerunOnFailure),
		)
	})
})
