package vmexport

import (
	"context"
	"fmt"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/ownerreference"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"go.uber.org/zap"

	kvcorev1 "kubevirt.io/api/core/v1"
	v1beta1 "kubevirt.io/api/export/v1beta1"
	snapshotv1 "kubevirt.io/api/snapshot/v1beta1"
	kubecli "kubevirt.io/client-go/kubecli"
)

const (
	sourceVM         string = "vm"
	sourceVMSnapshot string = "vmsnapshot"
	sourcePVC        string = "pvc"
)

func CreateVirtualMachineExport(virtClient kubecli.KubevirtClient, exportSourceKind, exportSourceNamespace, baseExportSourceName, secretName string) (*v1beta1.VirtualMachineExport, error) {
	source, err := getExportSource(exportSourceKind, baseExportSourceName)
	if err != nil {
		return nil, err
	}

	v1VmExport := &v1beta1.VirtualMachineExport{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: baseExportSourceName + "-",
			Namespace:    exportSourceNamespace,
		},
		Spec: v1beta1.VirtualMachineExportSpec{
			TokenSecretRef: &secretName,
			Source:         source,
		},
	}

	if err := ownerreference.SetPodOwnerReference(virtClient, v1VmExport); err != nil {
		return nil, err
	}

	return virtClient.VirtualMachineExport(exportSourceNamespace).Create(context.Background(), v1VmExport, metav1.CreateOptions{})
}

func WaitUntilVirtualMachineExportReady(client kubecli.KubevirtClient, namespace, name string) error {
	pollInterval := 15 * time.Second
	pollTimeout := 3600 * time.Second
	poller := func(ctx context.Context) (bool, error) {
		vmExport, err := client.VirtualMachineExport(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if vmExport.Status != nil {
			log.Logger().Info("VirtualMachineExport object status", zap.String("status", string(vmExport.Status.Phase)))

			if vmExport.Status.Phase == v1beta1.Ready {
				log.Logger().Info("VirtualMachineExport is in Ready state, and export source is not longer used")
				return true, nil
			}

			if vmExport.Status.Phase == v1beta1.Pending {
				log.Logger().Info("VirtualMachineExport is in Pending state, and export source is used")
				return false, nil
			}
		}
		return false, nil
	}

	return wait.PollUntilContextTimeout(context.Background(), pollInterval, pollTimeout, true, poller)
}

func GetRawDiskUrlFromVolumes(client kubecli.KubevirtClient, namespace, name, volumeName string) (string, error) {
	vmExport, err := client.VirtualMachineExport(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if vmExport.Status.Links == nil || vmExport.Status.Links.Internal == nil {
		return "", fmt.Errorf("no links found in VirtualMachineExport status")
	}

	for _, volume := range vmExport.Status.Links.Internal.Volumes {
		if volumeName != volume.Name {
			continue
		}

		for _, format := range volume.Formats {
			if format.Format == v1beta1.KubeVirtRaw {
				return format.Url, nil
			}
		}
	}
	return "", fmt.Errorf("volume %s is not found in VirtualMachineExport internal volumes", volumeName)
}

func GetLabelsFromExportSource(virtClient kubecli.KubevirtClient, exportSourceKind, exportSourceNamespace, exportSourceName, volumeName string) (map[string]string, error) {
	switch exportSourceKind {
	case sourceVM, sourceVMSnapshot:
		return getLabelsFromVirtualMachineOrSnapshot(virtClient, exportSourceNamespace, volumeName)
	case sourcePVC:
		return getLabelsFromPVC(virtClient, exportSourceNamespace, exportSourceName)
	default:
		return nil, fmt.Errorf("unsupported source kind: %s", exportSourceKind)
	}
}

func getLabelsFromVirtualMachineOrSnapshot(virtClient kubecli.KubevirtClient, namespace, volumeName string) (map[string]string, error) {
	dv, err := virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Get(context.Background(), volumeName, metav1.GetOptions{})
	if err == nil {
		return dv.GetLabels(), nil
	}
	if !errors.IsNotFound(err) {
		return nil, err
	}
	return getLabelsFromPVC(virtClient, namespace, volumeName)
}

func getLabelsFromPVC(virtClient kubecli.KubevirtClient, namespace, name string) (map[string]string, error) {
	pvc, err := virtClient.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pvc.GetLabels(), nil
}

func getExportSource(exportSourceKind, exportSourceName string) (corev1.TypedLocalObjectReference, error) {
	switch exportSourceKind {
	case sourceVM:
		return corev1.TypedLocalObjectReference{
			APIGroup: &kvcorev1.SchemeGroupVersion.Group,
			Kind:     "VirtualMachine",
			Name:     exportSourceName,
		}, nil
	case sourceVMSnapshot:
		return corev1.TypedLocalObjectReference{
			APIGroup: &snapshotv1.SchemeGroupVersion.Group,
			Kind:     "VirtualMachineSnapshot",
			Name:     exportSourceName,
		}, nil
	case sourcePVC:
		return corev1.TypedLocalObjectReference{
			APIGroup: &corev1.SchemeGroupVersion.Group,
			Kind:     "PersistentVolumeClaim",
			Name:     exportSourceName,
		}, nil
	default:
		return corev1.TypedLocalObjectReference{}, fmt.Errorf("invalid export-source-kind: %s, must be one of vm, vmsnapshot, pvc", exportSourceKind)
	}
}
