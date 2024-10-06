package certificate

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubecli "kubevirt.io/client-go/kubecli"
)

func GetCertificateFromVirtualMachineExport(client kubecli.KubevirtClient, namespace, name string) (string, error) {
	vmExport, err := client.VirtualMachineExport(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if vmExport.Status.Links == nil && vmExport.Status.Links.Internal == nil {
		return "", fmt.Errorf("no links found in VirtualMachineExport status")
	}

	content := vmExport.Status.Links.Internal.Cert
	if content == "" {
		return "", fmt.Errorf("no certificate found in VirtualMachineExport status")
	}
	return content, nil
}

func CreateCertificateFile(path, data string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %w", err)
	}
	return nil
}
