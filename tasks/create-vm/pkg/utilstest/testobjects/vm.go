package testobjects

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/rand"
	v1 "kubevirt.io/client-go/api/v1"
)

const NamespaceTestDefault = "default"

func newRandomVMI() *v1.VirtualMachineInstance {
	return newRandomVMIWithNS(NamespaceTestDefault)
}

func newRandomVMIWithNS(namespace string) *v1.VirtualMachineInstance {
	vmi := v1.NewMinimalVMIWithNS(namespace, "testvmi"+rand.String(48))

	vmi.Spec.Domain.Devices = v1.Devices{Interfaces: []v1.Interface{{Name: "default",
		InterfaceBindingMethod: v1.InterfaceBindingMethod{
			Masquerade: &v1.InterfaceMasquerade{}}}}}

	vmi.Spec.Networks = []v1.Network{*v1.DefaultPodNetwork()}

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

func NewTestVM() *v1.VirtualMachine {
	return newRandomVirtualMachine(newRandomVMI(), false)
}
