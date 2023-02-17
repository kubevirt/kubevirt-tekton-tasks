package test

import (
	"context"
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datavolume"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/dataobject"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/testconfigs"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/vm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

var _ = Describe("Create VM", func() {
	f := framework.NewFramework()
	BeforeEach(func() {
		if f.TestOptions.SkipCreateVMFromManifestTests {
			Skip("skipCreateVMFromManifestTests is set to true, skipping tests")
		}
	})

	ownedDVTestCase := map[datavolume.TestDataVolumeAttachmentType]int{datavolume.OwnedDV: 1}
	pvcTestCase := map[datavolume.TestDataVolumeAttachmentType]int{datavolume.PVC: 1}
	ownedPVCTestCase := map[datavolume.TestDataVolumeAttachmentType]int{datavolume.OwnedPVC: 2}
	dvTestCase := map[datavolume.TestDataVolumeAttachmentType]int{datavolume.DV: 2}
	allTestCase := map[datavolume.TestDataVolumeAttachmentType]int{datavolume.OwnedDV: 2, datavolume.OwnedPVC: 1, datavolume.PVC: 1, datavolume.DV: 1}

	DescribeTable("VM with attached PVCs/DV is created successfully", func(createMode CreateVMMode, idx int, runConf map[datavolume.TestDataVolumeAttachmentType]int) {
		mode := "template-mode"

		if createMode == CreateVMVMManifestMode {
			mode = "manifest-mode"
		}
		var datavolumes []*datavolume.TestDataVolume
		for attachmentType, count := range runConf {
			for id := 0; id < count; id++ {
				datavolumes = append(datavolumes,
					datavolume.NewBlankDataVolume(fmt.Sprintf("attach-to-vm-%s-%v-%v", mode, attachmentType, id)).AttachAs(attachmentType),
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

		for _, dv := range config.TaskData.DataVolumesToCreate {
			// wait for each DV to finish import, otherwise test will fail, because of not finished import of DV
			err := dataobject.WaitForSuccessfulDataVolume(f.KubevirtClient, dv.Data.Namespace, dv.Data.Name, constants.Timeouts.SmallDVCreation.Duration)
			Expect(err).ShouldNot(HaveOccurred())
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
		// check all disks are present
		Expect(vm.Spec.Template.Spec.Volumes).To(ConsistOf(expectedVM.Spec.Template.Spec.Volumes))
		Expect(vm.Spec.Template.Spec.Domain.Devices.Disks).To(ConsistOf(expectedVM.Spec.Template.Spec.Domain.Devices.Disks))
	},
		Entry("owned-dv=1", CreateVMTemplateMode, 0, ownedDVTestCase),
		Entry("pvc=1", CreateVMTemplateMode, 1, pvcTestCase),
		Entry("owned-pvc=2", CreateVMTemplateMode, 2, ownedPVCTestCase),
		Entry("dv=2", CreateVMTemplateMode, 3, dvTestCase),
		Entry("owned-dv=2,owned-pvc=1,pvc=1,dv=1", CreateVMTemplateMode, 4, allTestCase),
		Entry("owned-dv=1", CreateVMVMManifestMode, 0, ownedDVTestCase),
		Entry("pvc=1", CreateVMVMManifestMode, 1, pvcTestCase),
		Entry("owned-pvc=2", CreateVMVMManifestMode, 2, ownedPVCTestCase),
		Entry("dv=2", CreateVMVMManifestMode, 3, dvTestCase),
		Entry("owned-dv=2,owned-pvc=1,pvc=1,dv=1", CreateVMVMManifestMode, 4, allTestCase),
	)
	DescribeTable("VM with attached PVCs/DVs and existing disks/volumes is created successfully", func(createMode CreateVMMode) {
		mode := "template-mode"

		if createMode == CreateVMVMManifestMode {
			mode = "manifest-mode"
		}

		datavolumes := []*datavolume.TestDataVolume{
			datavolume.NewBlankDataVolume("attach-to-vm-with-disk-name-1-" + mode).AttachWithDiskName("disk1").AttachAs(datavolume.OwnedPVC),
			datavolume.NewBlankDataVolume("attach-to-vm-with-disk-name-2-" + mode).AttachWithDiskName("disk2").AttachAs(datavolume.OwnedDV),
			datavolume.NewBlankDataVolume("attach-to-vm-with-disk-name-3-" + mode).AttachWithDiskName("disk3").AttachAs(datavolume.OwnedDV),
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
				PersistentVolumeClaim: &kubevirtv1.PersistentVolumeClaimVolumeSource{
					PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "other",
					},
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
		for _, dvWrapper := range config.TaskData.DataVolumesToCreate {
			dataVolume, err := f.CdiClient.DataVolumes(dvWrapper.Data.Namespace).Create(context.TODO(), dvWrapper.Data, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageDataVolumes(dataVolume)
			config.TaskData.SetDVorPVC(fmt.Sprintf("%v:%v", dvWrapper.DiskName, dataVolume.Name), dvWrapper.AttachmentType)
		}

		if template := config.TaskData.Template; template != nil {
			template, err := f.TemplateClient.Templates(template.Namespace).Create(context.TODO(), template, v1.CreateOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.ManageTemplates(template)
		}

		for _, dv := range config.TaskData.DataVolumesToCreate {
			// wait for each DV to finish import, otherwise test will fail, because of not finished import of DV
			err := dataobject.WaitForSuccessfulDataVolume(f.KubevirtClient, dv.Data.Namespace, dv.Data.Name, constants.Timeouts.SmallDVCreation.Duration)
			Expect(err).ShouldNot(HaveOccurred())
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
		// check all disks are present
		Expect(vm.Spec.Template.Spec.Volumes).To(ConsistOf(expectedVM.Spec.Template.Spec.Volumes))
		Expect(vm.Spec.Template.Spec.Domain.Devices.Disks).To(ConsistOf(expectedVM.Spec.Template.Spec.Domain.Devices.Disks))
		Expect(vm.Spec.Template.Spec.Domain.Devices.Disks[2].CDRom.Bus).To(Equal(kubevirtv1.DiskBusSATA))
	},
		Entry(string(CreateVMTemplateMode), CreateVMTemplateMode),
		Entry(string(CreateVMVMManifestMode), CreateVMVMManifestMode),
	)
})
