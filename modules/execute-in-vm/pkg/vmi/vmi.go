package vmi

import (
	"errors"

	v1 "kubevirt.io/api/core/v1"
)

func GetPodIPAddress(vmi *v1.VirtualMachineInstance) (string, error) {
	podNetworkName := ""
	for _, network := range vmi.Spec.Networks {
		if network.Pod != nil {
			podNetworkName = network.Name
			break
		}
	}
	if podNetworkName == "" {
		return "", errors.New("pod network not found")
	}

	for _, statusInterface := range vmi.Status.Interfaces {
		if statusInterface.Name == podNetworkName {
			return statusInterface.IP, nil
		}
	}
	return "", nil
}
