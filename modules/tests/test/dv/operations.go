package dv

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdicliv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"
)

func DeleteDataVolume(dataVolumeClient cdicliv1beta1.DataVolumeInterface, dvName string, debug bool) {
	if !debug {
		_ = dataVolumeClient.Delete(dvName, &metav1.DeleteOptions{})
	}
}
