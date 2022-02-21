package test

import (
	"context"
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datavolume"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/vm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Create VM", func() {
	f := framework.NewFramework()

	for _, c := range []CreateVMMode{CreateVMTemplateMode, CreateVMVMManifestMode} {
		createMode := c
		Context(string(createMode), func() {
			Describe("VM with attached PVCs/DV is created successfully", func() {
				runConfigurations := []map[datavolume.TestDataVolumeAttachmentType]int{
					{
						datavolume.OwnedDV: 1,
					},
					{
						datavolume.PVC: 1,
					},
					{
						datavolume.OwnedPVC: 2,
					},
					{
						datavolume.DV: 2,
					},
					{
						// try all at once
						datavolume.OwnedDV:  2,
						datavolume.OwnedPVC: 1,
						datavolume.PVC:      1,
						datavolume.DV:       1,
					},
				}

				for i, r := range runConfigurations {
					idx, runConf := i, r
					name := ""
					for attachmentType, count := range runConf {
						name += fmt.Sprintf("%v=%v ", attachmentType, count)
					}
					It(name, func() {
						var datavolumes []*datavolume.TestDataVolume
						for attachmentType, count := range runConf {
							name += fmt.Sprintf("%v=%v ", attachmentType, count)
							for id := 0; id < count; id++ {
								datavolumes = append(datavolumes,
									datavolume.NewBlankDataVolume(fmt.Sprintf("attach-to-vm-%v-%v", attachmentType, id)).AttachAs(attachmentType),
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
									Timeout:        Timeouts.SmallDVCreation,
								},
								TaskData: testconfigs.CreateVMTaskData{
									CreateMode:                createMode,
									VM:                        testobjects.NewTestAlpineVM("create-vm-from-manifest-with-disk").Build(),
									DataVolumesToCreate:       datavolumes,
									ExpectedAdditionalDiskBus: "virtio",
								},
							}
						case CreateVMTemplateMode:
							expectedDiskBus := "virtio"
							testTemplate := template.NewCirrosServerTinyTemplate()
							switch idx % 4 { // try different disk buses for each test
							case 0:
								testTemplate.WithSataDiskValidations()
								expectedDiskBus = "sata"
							case 1:
								testTemplate.WithSCSIDiskValidations()
								expectedDiskBus = "scsi"
							case 2:
								testTemplate.WithVirtioDiskValidations()
							}
							config = &testconfigs.CreateVMTestConfig{
								TaskRunTestConfig: testconfigs.TaskRunTestConfig{
									ServiceAccount: CreateVMFromTemplateServiceAccountName,
									ExpectedLogs:   ExpectedSuccessfulVMCreation,
									Timeout:        Timeouts.SmallDVCreation,
									LimitEnvScope:  OKDEnvScope,
								},
								TaskData: testconfigs.CreateVMTaskData{
									CreateMode: createMode,
									Template:   testTemplate.Build(),
									TemplateParams: []string{
										template.TemplateParam(template.NameParam, E2ETestsRandomName("create-vm-from-template-with-disk")),
									},
									DataVolumesToCreate:       datavolumes,
									ExpectedAdditionalDiskBus: expectedDiskBus,
								},
							}
						default:
							panic("invalid create mode")
						}

						f.TestSetup(config)
						if template := config.TaskData.Template; template != nil {
							template, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
							Expect(err).ShouldNot(HaveOccurred())
							f.ManageTemplates(template)
						}
						for _, dvWrapper := range config.TaskData.DataVolumesToCreate {
							dataVolume, err := f.CdiClient.DataVolumes(dvWrapper.Data.Namespace).Create(context.TODO(), dvWrapper.Data, v1.CreateOptions{})
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
		It("VM with attached PVCs/DVs and existing disks/volumes is created successfully", func() {
			datavolumes := []*datavolume.TestDataVolume{
				datavolume.NewBlankDataVolume("attach-to-vm-with-disk-name-1").AttachWithDiskName("disk1").AttachAs(datavolume.OwnedPVC),
				datavolume.NewBlankDataVolume("attach-to-vm-with-disk-name-2").AttachWithDiskName("disk2").AttachAs(datavolume.OwnedDV),
				datavolume.NewBlankDataVolume("attach-to-vm-with-disk-name-3").AttachWithDiskName("disk3").AttachAs(datavolume.OwnedDV),
			}

			vmDisk1 := kubevirtv1.Disk{
				Name: datavolumes[0].DiskName,
				DiskDevice: kubevirtv1.DiskDevice{
					CDRom: &kubevirtv1.CDRomTarget{Bus: "sata"},
				},
			}
			vmDisk2 := kubevirtv1.Disk{
				Name: datavolumes[1].DiskName,
				DiskDevice: kubevirtv1.DiskDevice{
					Disk: &kubevirtv1.DiskTarget{Bus: "virtio"},
				},
			}

			// disk disk3 should be created by the task

			// volume disk1 should be created by the task

			vmVolume2 := kubevirtv1.Volume{
				Name: datavolumes[1].DiskName,
				// wrong source - should overwrite
				VolumeSource: kubevirtv1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "other",
					},
				},
			}
			vmVolume3 := kubevirtv1.Volume{
				Name: datavolumes[2].DiskName,
				// no source - should complete
			}

			var config *testconfigs.CreateVMTestConfig

			switch createMode {
			case CreateVMVMManifestMode:
				config = &testconfigs.CreateVMTestConfig{
					TaskRunTestConfig: testconfigs.TaskRunTestConfig{
						ServiceAccount: CreateVMFromManifestServiceAccountName,
						ExpectedLogs:   ExpectedSuccessfulVMCreation,
						Timeout:        Timeouts.SmallDVCreation,
					},
					TaskData: testconfigs.CreateVMTaskData{
						CreateMode: createMode,
						VM: testobjects.NewTestAlpineVM("create-vm-from-manifest-with-existing-disk").
							// to be compatible with the template flow
							WithCloudConfig(
								testobjects.CloudConfig{
									Password: "alpine",
								},
							).
							WithDisk(vmDisk1).
							WithDisk(vmDisk2).
							WithVolume(vmVolume2).
							WithVolume(vmVolume3).
							Build(),
						DataVolumesToCreate:       datavolumes,
						ExpectedAdditionalDiskBus: "virtio",
					},
				}
			case CreateVMTemplateMode:
				config = &testconfigs.CreateVMTestConfig{
					TaskRunTestConfig: testconfigs.TaskRunTestConfig{
						ServiceAccount: CreateVMFromTemplateServiceAccountName,
						ExpectedLogs:   ExpectedSuccessfulVMCreation,
						Timeout:        Timeouts.SmallDVCreation,
						LimitEnvScope:  OKDEnvScope,
					},
					TaskData: testconfigs.CreateVMTaskData{
						CreateMode: createMode,
						Template: template.NewCirrosServerTinyTemplate().
							WithDisk(vmDisk1).
							WithDisk(vmDisk2).
							WithVolume(vmVolume2).
							WithVolume(vmVolume3).
							Build(),
						TemplateParams: []string{
							template.TemplateParam(template.NameParam, E2ETestsRandomName("create-vm-from-template-with-existing-disk")),
						},
						DataVolumesToCreate:       datavolumes,
						ExpectedAdditionalDiskBus: "virtio",
					},
				}
			default:
				panic("invalid create mode")
			}

			f.TestSetup(config)
			if template := config.TaskData.Template; template != nil {
				template, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageTemplates(template)
			}
			for _, dvWrapper := range config.TaskData.DataVolumesToCreate {
				dataVolume, err := f.CdiClient.DataVolumes(dvWrapper.Data.Namespace).Create(context.TODO(), dvWrapper.Data, v1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				f.ManageDataVolumes(dataVolume)
				config.TaskData.SetDVorPVC(fmt.Sprintf("%v:%v", dvWrapper.DiskName, dataVolume.Name), dvWrapper.AttachmentType)
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
			Expect(vm.Spec.Template.Spec.Domain.Devices.Disks[2].CDRom.Bus).To(Equal("sata"))
		})
	}
})
