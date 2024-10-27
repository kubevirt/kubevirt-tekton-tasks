package certificate

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubecli "kubevirt.io/client-go/kubecli"
)

func GetCertificateFromVirtualMachineExport(client kubecli.KubevirtClient, namespace, name string) (string, error) {
	vmExport, err := client.VirtualMachineExport(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if vmExport.Status.Links == nil || vmExport.Status.Links.Internal == nil {
		return "", fmt.Errorf("no links found in VirtualMachineExport status")
	}

	content := vmExport.Status.Links.Internal.Cert
	if content == "" {
		return "", fmt.Errorf("no certificate found in VirtualMachineExport status")
	}
	return content, nil
}
