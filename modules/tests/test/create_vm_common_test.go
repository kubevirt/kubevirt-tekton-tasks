package test

import (
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/dv"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	templ "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/template"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/vm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create VM", func() {
	f := framework.NewFramework()

	for _, createMode := range []CreateVMMode{CreateVMTemplateMode, CreateVMVMManifestMode} {
		Context(string(createMode), func() {
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

						var config *testconfigs.CreateVMTestConfig

						switch createMode {
						case CreateVMVMManifestMode:
							config = &testconfigs.CreateVMTestConfig{
								TaskRunTestConfig: testconfigs.TaskRunTestConfig{
									ServiceAccount: CreateVMFromManifestServiceAccountName,
									ExpectedLogs:   ExpectedSuccessfulVMCreation,
									Timeout:        Timeouts.SmallBlankDVCreation,
								},
								TaskData: testconfigs.CreateVMTaskData{
									CreateMode:                createMode,
									VM:                        testobjects.NewTestAlpineVM("create-vm-from-manifest-with-disk").Build(),
									DataVolumesToCreate:       datavolumes,
									ExpectedAdditionalDiskBus: "virtio",
								},
							}
						case CreateVMTemplateMode:
							expectedDisbBus := "virtio"
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
							}
							config = &testconfigs.CreateVMTestConfig{
								TaskRunTestConfig: testconfigs.TaskRunTestConfig{
									ServiceAccount: CreateVMFromTemplateServiceAccountName,
									ExpectedLogs:   ExpectedSuccessfulVMCreation,
									Timeout:        Timeouts.SmallBlankDVCreation,
									LimitEnvScope:  OpenshiftEnvScope,
								},
								TaskData: testconfigs.CreateVMTaskData{
									CreateMode: createMode,
									Template:   testTemplate.Build(),
									TemplateParams: []string{
										templ.TemplateParam(templ.NameParam, E2ETestsRandomName("create-vm-from-template-with-disk")),
									},
									DataVolumesToCreate:       datavolumes,
									ExpectedAdditionalDiskBus: expectedDisbBus,
								},
							}
						default:
							panic("invalid create mode")
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
							ExpectLogs(config.GetAllExpectedLogs()...).
							ExpectResults(map[string]string{
								CreateVMResults.Name:      expectedVM.Name,
								CreateVMResults.Namespace: expectedVM.Namespace,
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
		})
	}
})
