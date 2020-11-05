package dv

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	cdiv1beta12 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
	cdicliv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"
	"time"
)

func WaitForSuccessfulDataVolume(cdiClientSet cdicliv1beta1.CdiV1beta1Interface, namespace, name string, timeout time.Duration) error {
	return wait.PollImmediate(constants.PollInterval, timeout, func() (bool, error) {
		dataVolume, err := cdiClientSet.DataVolumes(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return isDataVolumeImportStatusSuccessful(dataVolume), nil
	})
}

func IsDataVolumeImportSuccessful(cdiClientSet cdicliv1beta1.CdiV1beta1Interface, namespace string, name string) bool {
	dataVolume, err := cdiClientSet.DataVolumes(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return isDataVolumeImportStatusSuccessful(dataVolume)
}

func isDataVolumeImportStatusSuccessful(dataVolume *cdiv1beta12.DataVolume) bool {
	return GetConditionMap(dataVolume)[cdiv1beta12.DataVolumeBound] == v1.ConditionTrue &&
		dataVolume.Status.Phase == cdiv1beta12.Succeeded
}
