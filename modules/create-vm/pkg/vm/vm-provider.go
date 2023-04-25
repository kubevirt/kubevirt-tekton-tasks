package vm

import (
	"context"

	kubevirtv1 "kubevirt.io/api/core/v1"
	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
)

type virtualMachineProvider struct {
	client kubevirtcliv1.KubevirtClient
}

type VirtualMachineProvider interface {
	Create(namespace string, vm *kubevirtv1.VirtualMachine) (*kubevirtv1.VirtualMachine, error)
	Start(namespace, name string) error
}

func NewVirtualMachineProvider(client kubevirtcliv1.KubevirtClient) VirtualMachineProvider {
	return &virtualMachineProvider{
		client: client,
	}
}

func (v *virtualMachineProvider) Create(namespace string, vm *kubevirtv1.VirtualMachine) (*kubevirtv1.VirtualMachine, error) {
	return v.client.VirtualMachine(namespace).Create(context.Background(), vm)
}

func (v *virtualMachineProvider) Start(namespace, name string) error {
	return v.client.VirtualMachine(namespace).Start(context.Background(), name, &kubevirtv1.StartOptions{})
}
