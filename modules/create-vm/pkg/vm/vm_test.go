package vm_test

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utilstest/testobjects"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	"sort"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates/validations"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	vm2 "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/vm"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	shtestobjects "github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects"
)

var _ = Describe("VM", func() {
	var vm *kubevirtv1.VirtualMachine

	BeforeEach(func() {
		vm = shtestobjects.NewTestVM()
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
			vm = shtestobjects.NewTestVM()
			emptyValidations = validations.NewTemplateValidations(nil)
			cliOptions = &parse.CLIOptions{
				TemplateName:              "test",
				TemplateNamespaces:        []string{"default"},
				VirtualMachineNamespaces:  []string{"default"},
				OwnDataVolumes:            []string{"dv1"},
				DataVolumes:               []string{"dv2", "dv3"},
				OwnPersistentVolumeClaims: []string{"pvc1", "pvc2"},
				PersistentVolumeClaims:    []string{"pvc3"},
			}
			Expect(cliOptions.Init()).Should(Succeed())
		})

		table.DescribeTable("adds all volumes with various validations", func(templateValidations *validations.TemplateValidations, expectedBus string) {
			addsVolumesCorrectly(vm, templateValidations, cliOptions, expectedBus)
		},
			table.Entry("no validations", nil, Virtio),
			table.Entry("empty validations", validations.NewTemplateValidations(nil), Virtio),
			table.Entry("empty validations", validations.NewTemplateValidations(testobjects.NewTestCommonTemplateValidations(Scsi)), Scsi),
		)

		It("adds some volumes", func() {
			cliOptions.DataVolumes = nil
			cliOptions.OwnPersistentVolumeClaims = nil
			addsVolumesCorrectly(vm, emptyValidations, cliOptions, Virtio)
		})

		It("adds no volumes", func() {
			cliOptions.OwnDataVolumes = nil
			cliOptions.DataVolumes = nil
			cliOptions.OwnPersistentVolumeClaims = nil
			cliOptions.PersistentVolumeClaims = nil
			addsVolumesCorrectly(vm, emptyValidations, cliOptions, Virtio)
		})
	})

	It("Adds correct metadata", func() {
		vm2.AddMetadata(vm, shtestobjects.NewFedoraServerTinyTemplate())

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
})

func addsVolumesCorrectly(vm *kubevirtv1.VirtualMachine, templateValidations *validations.TemplateValidations, cliOpts *parse.CLIOptions, expectedBus string) {
	vm2.AddVolumes(vm, templateValidations, cliOpts)
	disksCount := len(cliOpts.GetAllDiskNames())
	Expect(vm.Spec.Template.Spec.Volumes).To(HaveLen(disksCount))
	Expect(vm.Spec.Template.Spec.Domain.Devices.Disks).To(HaveLen(disksCount))

	var foundDiskNames []string
	var foundVolumeNames []string
	for _, disk := range vm.Spec.Template.Spec.Domain.Devices.Disks {
		Expect(disk.Disk.Bus).To(Equal(expectedBus))
		foundDiskNames = append(foundDiskNames, disk.Name)
	}
	for _, volume := range vm.Spec.Template.Spec.Volumes {
		foundVolumeNames = append(foundVolumeNames, volume.Name)
	}

	var expectedNames = cliOpts.GetAllDiskNames()
	sort.Strings(foundDiskNames)
	sort.Strings(foundVolumeNames)
	sort.Strings(expectedNames)

	Expect(foundDiskNames).To(Equal(expectedNames))
	Expect(foundVolumeNames).To(Equal(expectedNames))
}
