package vm_test

import (
	"sort"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest/testobjects"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates/validations"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	vm2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/vm"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	shtestobjects "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
	template "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/template"
)

var _ = Describe("VM", func() {
	var vm *kubevirtv1.VirtualMachine

	BeforeEach(func() {
		vm = shtestobjects.NewTestVM().Build()
	})

	It("Creates owner reference", func() {
		result := vm2.AsVMOwnerReference(vm)
		Expect(result).To(BeAssignableToTypeOf(metav1.OwnerReference{}))
		Expect(result.Name).To(Equal(vm.Name))
	})

	Describe("Adds volumes", func() {
		var emptyValidations *validations.TemplateValidations
		var cliOptions *parse.CLIOptions

		BeforeEach(func() {
			vm = shtestobjects.NewTestVM().Build()
			emptyValidations = validations.NewTemplateValidations(nil)
			cliOptions = &parse.CLIOptions{
				TemplateName:              "test",
				TemplateNamespace:         "default",
				VirtualMachineNamespace:   "default",
				PersistentVolumeClaims:    []string{"pvc1"},
				OwnPersistentVolumeClaims: []string{"pvc2", "pvc3"},
				DataVolumes:               []string{"dv1", "dv2"},
				OwnDataVolumes:            []string{"dv3"},
			}
			Expect(cliOptions.Init()).Should(Succeed())
		})

		DescribeTable("adds all volumes with various validations", func(templateValidations *validations.TemplateValidations, expectedBus string) {
			addsVolumesCorrectly(vm, templateValidations, cliOptions, []string{expectedBus})
		},
			Entry("no validations", nil, Virtio),
			Entry("empty validations", validations.NewTemplateValidations(nil), Virtio),
			Entry("empty validations", validations.NewTemplateValidations(testobjects.NewTestCommonTemplateValidations(Scsi)), Scsi),
		)

		It("adds some volumes", func() {
			cliOptions.DataVolumes = nil
			cliOptions.OwnPersistentVolumeClaims = nil
			addsVolumesCorrectly(vm, emptyValidations, cliOptions, []string{Virtio})
		})

		It("adds no volumes", func() {
			cliOptions.OwnDataVolumes = nil
			cliOptions.DataVolumes = nil
			cliOptions.OwnPersistentVolumeClaims = nil
			cliOptions.PersistentVolumeClaims = nil
			addsVolumesCorrectly(vm, emptyValidations, cliOptions, []string{Virtio})
		})

		It("adding named disks", func() {
			cliOptions.OwnDataVolumes = append(cliOptions.OwnDataVolumes, "disk1:dv4")
			cliOptions.PersistentVolumeClaims = append(cliOptions.PersistentVolumeClaims, "disk2:pvc4")
			addsVolumesCorrectly(vm, emptyValidations, cliOptions, []string{Virtio})
		})

		It("replaces existing disks", func() {
			bootOrder := uint(1)
			vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks,
				kubevirtv1.Disk{
					Name:      "disk1",
					BootOrder: &bootOrder,
					DiskDevice: kubevirtv1.DiskDevice{
						CDRom: &kubevirtv1.CDRomTarget{Bus: Sata},
					},
				},
				kubevirtv1.Disk{
					Name: "disk2",
					DiskDevice: kubevirtv1.DiskDevice{
						Disk: &kubevirtv1.DiskTarget{Bus: Virtio},
					},
				},
				kubevirtv1.Disk{
					Name: "disk3",
					DiskDevice: kubevirtv1.DiskDevice{
						Disk: &kubevirtv1.DiskTarget{Bus: Sata},
					},
				},
			)
			vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes,
				kubevirtv1.Volume{
					Name: "disk1",
					// wrong source - should overwrite
					VolumeSource: kubevirtv1.VolumeSource{
						PersistentVolumeClaim: &kubevirtv1.PersistentVolumeClaimVolumeSource{
							PersistentVolumeClaimVolumeSource: v1.PersistentVolumeClaimVolumeSource{
								ClaimName: "other1",
							},
						},
					},
				},
				kubevirtv1.Volume{
					// no source - should complete
					Name: "disk2",
				},
				// for disk3 - should create volume
				// for disk4 - should not damage source
				kubevirtv1.Volume{
					Name: "disk4",
					VolumeSource: kubevirtv1.VolumeSource{
						PersistentVolumeClaim: &kubevirtv1.PersistentVolumeClaimVolumeSource{
							PersistentVolumeClaimVolumeSource: v1.PersistentVolumeClaimVolumeSource{
								ClaimName: "other2",
								ReadOnly:  true,
							},
						},
					},
				},
			)

			cliOptions.OwnDataVolumes = append(cliOptions.OwnDataVolumes, "disk1:dv4")
			cliOptions.PersistentVolumeClaims = append(cliOptions.PersistentVolumeClaims, "disk2:pvc4", "disk3:pvc5", "disk4:pvc6")
			addsVolumesCorrectly(vm, emptyValidations, cliOptions, []string{Sata, Virtio, Sata, Virtio})
			// initial disks should not be changed
			Expect(vm.Spec.Template.Spec.Domain.Devices.Disks[0]).Should(Equal(
				kubevirtv1.Disk{
					Name:      "disk1",
					BootOrder: &bootOrder,
					DiskDevice: kubevirtv1.DiskDevice{
						CDRom: &kubevirtv1.CDRomTarget{Bus: Sata},
					},
				},
			))
			Expect(vm.Spec.Template.Spec.Domain.Devices.Disks[1]).Should(Equal(
				kubevirtv1.Disk{
					Name: "disk2",
					DiskDevice: kubevirtv1.DiskDevice{
						Disk: &kubevirtv1.DiskTarget{Bus: Virtio},
					},
				},
			))
			Expect(vm.Spec.Template.Spec.Domain.Devices.Disks[2]).Should(Equal(
				kubevirtv1.Disk{
					Name: "disk3",
					DiskDevice: kubevirtv1.DiskDevice{
						Disk: &kubevirtv1.DiskTarget{Bus: Sata},
					},
				},
			))
			// should have correctly filled sources
			Expect(vm.Spec.Template.Spec.Volumes[0]).Should(Equal(
				kubevirtv1.Volume{
					Name: "disk1",
					VolumeSource: kubevirtv1.VolumeSource{
						DataVolume: &kubevirtv1.DataVolumeSource{Name: "dv4"},
					},
				},
			))
			Expect(vm.Spec.Template.Spec.Volumes[1]).Should(Equal(
				kubevirtv1.Volume{
					Name: "disk2",
					VolumeSource: kubevirtv1.VolumeSource{
						PersistentVolumeClaim: &kubevirtv1.PersistentVolumeClaimVolumeSource{
							PersistentVolumeClaimVolumeSource: v1.PersistentVolumeClaimVolumeSource{
								ClaimName: "pvc4",
							},
						},
					},
				},
			))
			Expect(vm.Spec.Template.Spec.Volumes[2]).Should(Equal(
				kubevirtv1.Volume{
					Name: "disk4",
					VolumeSource: kubevirtv1.VolumeSource{
						PersistentVolumeClaim: &kubevirtv1.PersistentVolumeClaimVolumeSource{
							PersistentVolumeClaimVolumeSource: v1.PersistentVolumeClaimVolumeSource{
								ClaimName: "pvc6",
								ReadOnly:  true,
							},
						},
					},
				},
			))
		})
	})

	It("Adds correct metadata from template", func() {
		vm2.AddMetadata(vm, template.NewFedoraServerTinyTemplate().Build())

		Expect(vm.Labels).To(Equal(map[string]string{
			"vm.kubevirt.io/template":              "fedora-server-tiny-v0.7.0",
			"vm.kubevirt.io/template.namespace":    "openshift",
			"os.template.kubevirt.io/fedora29":     "true",
			"flavor.template.kubevirt.io/tiny":     "true",
			"workload.template.kubevirt.io/server": "true",
		}))

		Expect(vm.Annotations).To(Equal(map[string]string{
			"name.os.template.kubevirt.io/fedora29": "Fedora 27 or higher",
		}))

		Expect(vm.Spec.Template.ObjectMeta.Labels).To(Equal(map[string]string{
			"vm.kubevirt.io/name":                  vm.Name,
			"name":                                 vm.Name,
			"os.template.kubevirt.io/fedora29":     "true",
			"flavor.template.kubevirt.io/tiny":     "true",
			"workload.template.kubevirt.io/server": "true",
		}))

	})

	It("Adds correct default metadata", func() {
		vm2.AddMetadata(vm, nil)

		Expect(vm.Spec.Template.ObjectMeta.Labels).To(Equal(map[string]string{
			"vm.kubevirt.io/name": vm.Name,
			"name":                vm.Name,
		}))

	})
})

// expectedBuses: the disk at index i should have a bus at expectedBuses[i], or the last of expectedBuses
func addsVolumesCorrectly(vm *kubevirtv1.VirtualMachine, templateValidations *validations.TemplateValidations, cliOpts *parse.CLIOptions, expectedBuses []string) {
	vm2.AddVolumes(vm, templateValidations, cliOpts)
	disksCount := len(cliOpts.GetPVCDiskNamesMap()) + len(cliOpts.GetDVDiskNamesMap())
	Expect(vm.Spec.Template.Spec.Volumes).To(HaveLen(disksCount))
	Expect(vm.Spec.Template.Spec.Domain.Devices.Disks).To(HaveLen(disksCount))

	var foundDiskNames []string
	var foundVolumeNames []string
	for i, disk := range vm.Spec.Template.Spec.Domain.Devices.Disks {
		var expectedBus string
		if i < len(expectedBuses) {
			expectedBus = expectedBuses[i]
		} else {
			expectedBus = expectedBuses[len(expectedBuses)-1]
		}

		if disk.CDRom != nil {
			Expect(disk.CDRom.Bus).To(Equal(expectedBus))
		} else {
			Expect(disk.Disk.Bus).To(Equal(expectedBus))
		}

		foundDiskNames = append(foundDiskNames, disk.Name)
	}
	for _, volume := range vm.Spec.Template.Spec.Volumes {
		foundVolumeNames = append(foundVolumeNames, volume.Name)
	}

	var expectedNames []string

	for expectedName, _ := range cliOpts.GetPVCDiskNamesMap() {
		expectedNames = append(expectedNames, expectedName)
	}

	for expectedName, _ := range cliOpts.GetDVDiskNamesMap() {
		expectedNames = append(expectedNames, expectedName)
	}

	sort.Strings(foundDiskNames)
	sort.Strings(foundVolumeNames)
	sort.Strings(expectedNames)

	Expect(foundDiskNames).To(Equal(expectedNames))
	Expect(foundVolumeNames).To(Equal(expectedNames))
}
