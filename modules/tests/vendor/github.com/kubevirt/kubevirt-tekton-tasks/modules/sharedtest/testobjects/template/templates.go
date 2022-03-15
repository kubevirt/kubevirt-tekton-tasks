package template

import (
	"encoding/json"
	"fmt"

	v1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

const (
	descriptionAnnotation    = "description"
	validationsAnnotation    = "validations"
	DescriptionParam         = "DESCRIPTION"
	NameParam                = "NAME"
	DataVolumeNameParam      = "DATA_SOURCE_NAME"
	DataVolumeNamespaceParam = "DATA_SOURCE_NAMESPACE"
)

type TestTemplate struct {
	Data *v1.Template
}

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(kubevirtv1.AddToScheme(scheme))
}

func (t *TestTemplate) Build() *v1.Template {
	return t.Data
}

func (t *TestTemplate) modifyVM(processVM func(vm *kubevirtv1.VirtualMachine)) {
	for idx, obj := range t.Data.Objects {
		decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(kubevirtv1.GroupVersion)
		decoded, err := runtime.Decode(decoder, obj.Raw)
		if err != nil {
			panic(err)
		}
		vm, ok := decoded.(*kubevirtv1.VirtualMachine)
		if ok {
			processVM(vm)
			bytes, err := json.Marshal(vm)
			if err != nil {
				panic(err)
			}
			t.Data.Objects[idx] = runtime.RawExtension{
				Raw: bytes,
			}
			return
		}
	}
	panic("no vm found")

}

func (t *TestTemplate) WithDescriptionParam() *TestTemplate {
	t.modifyVM(func(vm *kubevirtv1.VirtualMachine) {
		if vm.Annotations == nil {
			vm.Annotations = map[string]string{}
		}
		vm.Annotations[descriptionAnnotation] = fmt.Sprintf("${%v}", DescriptionParam)
	})
	t.Data.Parameters = append(t.Data.Parameters, v1.Parameter{
		Name:        DescriptionParam,
		Description: DescriptionParam,
		Value:       "",
		Required:    true,
	})

	return t
}

func (t *TestTemplate) WithDisk(disk kubevirtv1.Disk) *TestTemplate {
	t.modifyVM(func(vm *kubevirtv1.VirtualMachine) {
		vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, disk)
	})

	return t
}

func (t *TestTemplate) WithVolume(volume kubevirtv1.Volume) *TestTemplate {
	t.modifyVM(func(vm *kubevirtv1.VirtualMachine) {
		vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, volume)
	})

	return t
}

func (t *TestTemplate) WithNonMatchingDisk() *TestTemplate {
	t.modifyVM(func(vm *kubevirtv1.VirtualMachine) {
		vm.Spec.Template.Spec.Domain.Devices.Disks[0].Name = "non-matching-name"
	})
	return t
}

func (t *TestTemplate) WithNoVM() *TestTemplate {
	t.Data.Objects = nil
	return t
}

func (t *TestTemplate) WithSCSIDiskValidations() *TestTemplate {
	t.Data.Annotations[validationsAnnotation] = `[
  {
    "name": "scsi-bus",
    "path": "jsonpath::.spec.domain.devices.disks[*].disk.bus",
    "rule": "enum",
    "values": ["scsi"],
    "justWarning": true
  }, {
    "name": "disk-bus",
    "path": "jsonpath::.spec.domain.devices.disks[*].disk.bus",
    "rule": "enum",
    "values": ["virtio", "scsi"]
  }
]`
	return t
}

func (t *TestTemplate) WithSataDiskValidations() *TestTemplate {
	t.Data.Annotations[validationsAnnotation] = `[
  {
    "name": "disk-bus",
    "path": "jsonpath::.spec.domain.devices.disks[*].disk.bus",
    "rule": "enum",
    "values": ["sata"]
  }
]`
	return t
}

func (t *TestTemplate) WithVirtioDiskValidations() *TestTemplate {
	t.Data.Annotations[validationsAnnotation] = `[
  {
    "name": "disk-bus",
    "path": "jsonpath::.spec.domain.devices.disks[*].disk.bus",
    "rule": "enum",
    "values": ["sata", "virtio"]
  }
]`
	return t
}
