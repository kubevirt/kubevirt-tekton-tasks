package testobjects

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/rand"
	v1 "kubevirt.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type TestVM struct {
	Data *v1.VirtualMachine
}

func newRandomVMI() *v1.VirtualMachineInstance {
	return newRandomVMIWithNS(testconstants.NamespaceTestDefault)
}

func newRandomVMIWithNS(namespace string) *v1.VirtualMachineInstance {
	vmi := v1.NewVMIReferenceFromNameWithNS(namespace, "testvmi"+rand.String(48))
	vmi.Spec.Domain.CPU = &v1.CPU{
		Cores:   1,
		Sockets: 1,
		Threads: 1,
	}
	vmi.Spec.Domain.Devices = v1.Devices{
		Interfaces: []v1.Interface{
			*v1.DefaultBridgeNetworkInterface(),
		},
	}

	vmi.Spec.Networks = []v1.Network{
		*v1.DefaultPodNetwork(),
	}

	return vmi
}

func newRandomVirtualMachine(vmi *v1.VirtualMachineInstance, running bool) *v1.VirtualMachine {
	name := vmi.Name
	namespace := vmi.Namespace
	labels := map[string]string{"name": name}
	for k, v := range vmi.Labels {
		labels[k] = v
	}
	vm := &v1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.VirtualMachineSpec{
			Running: &running,
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    labels,
					Name:      name + "makeitinteresting", // this name should have no effect
					Namespace: namespace,
				},
				Spec: vmi.Spec,
			},
		},
	}
	vm.SetGroupVersionKind(schema.GroupVersionKind{Group: v1.GroupVersion.Group, Kind: "VirtualMachine", Version: v1.GroupVersion.Version})
	return vm
}

func NewTestVM() *TestVM {
	return &TestVM{
		Data: newRandomVirtualMachine(newRandomVMI(), false),
	}
}

func NewTestVMI() *v1.VirtualMachineInstance {
	return newRandomVMI()
}

func (t *TestVM) Build() *v1.VirtualMachine {
	return t.Data
}

func (t *TestVM) ToString() string {
	outBytes, _ := yaml.Marshal(t.Data)
	return string(outBytes)
}

func (t *TestVM) WithMemory(memory string) *TestVM {
	t.Data.Spec.Template.Spec.Domain.Resources.Requests = corev1.ResourceList{
		corev1.ResourceMemory: resource.MustParse(memory),
	}
	return t
}

func (t *TestVM) WithDisk(disk v1.Disk) *TestVM {
	t.Data.Spec.Template.Spec.Domain.Devices.Disks = append(t.Data.Spec.Template.Spec.Domain.Devices.Disks, disk)
	return t
}

func (t *TestVM) WithVolume(volume v1.Volume) *TestVM {
	t.Data.Spec.Template.Spec.Volumes = append(t.Data.Spec.Template.Spec.Volumes, volume)
	return t
}

func (t *TestVM) WithNonMatchingDisk() *TestVM {
	t.Data.Spec.Template.Spec.Domain.Devices.Disks[0].Name = "non-matching-name"
	return t
}

func (t *TestVM) WithLabel(key, value string) *TestVM {
	if t.Data.Labels == nil {
		t.Data.Labels = map[string]string{}
	}
	t.Data.Labels[key] = value
	return t
}

func (t *TestVM) WithVMILabel(key, value string) *TestVM {
	if t.Data.Labels == nil {
		t.Data.Spec.Template.ObjectMeta.Labels = map[string]string{}
	}
	t.Data.Spec.Template.ObjectMeta.Labels[key] = value
	return t
}

func (t *TestVM) WithCloudConfig(cloudConfig CloudConfig) *TestVM {
	if cloudConfig.Password == "" {
		cloudConfig.Password = "fedora"
	}

	applied := false

	for _, volume := range t.Data.Spec.Template.Spec.Volumes {
		if volume.CloudInitNoCloud != nil {
			volume.CloudInitNoCloud.UserData = cloudConfig.ToString()
			volume.CloudInitNoCloud.UserDataBase64 = ""
			applied = true
			break
		}
	}

	if !applied {
		cloudinitDiskName := "cloudinitdisk"
		t.Data.Spec.Template.Spec.Domain.Devices.Disks = append(t.Data.Spec.Template.Spec.Domain.Devices.Disks, v1.Disk{
			Name: cloudinitDiskName,
			DiskDevice: v1.DiskDevice{
				Disk: &v1.DiskTarget{
					Bus: "virtio",
				},
			},
		})
		t.Data.Spec.Template.Spec.Volumes = append(t.Data.Spec.Template.Spec.Volumes, v1.Volume{
			Name: cloudinitDiskName,
			VolumeSource: v1.VolumeSource{
				CloudInitNoCloud: &v1.CloudInitNoCloudSource{
					UserData: cloudConfig.ToString(),
				},
			},
		})
	}
	return t
}
