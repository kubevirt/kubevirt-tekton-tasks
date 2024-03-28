package dataobject

import (
	"context"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	cdicliv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"
)

func WaitForSuccessfulDataVolume(kubevirtClient kubevirtcliv1.KubevirtClient, namespace, name string, timeout time.Duration) error {
	return wait.PollImmediate(constants.PollInterval, timeout, func() (bool, error) {
		return IsDataVolumeImportSuccessful(kubevirtClient, namespace, name)
	})
}

func WaitForSuccessfulDataSource(cdiClientSet cdicliv1beta1.CdiV1beta1Interface, namespace, name string, timeout time.Duration) error {
	return wait.PollImmediate(constants.PollInterval, timeout, func() (bool, error) {
		dataSource, err := cdiClientSet.DataSources(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return IsDataSourceReady(dataSource), nil
	})
}

func IsDataVolumeImportSuccessful(kubevirtClient kubevirtcliv1.KubevirtClient, namespace, name string) (bool, error) {
	dataVolume, err := kubevirtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		pvc, err := kubevirtClient.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if pvc != nil {
			return true, nil
		}
		return true, err
	} else if err != nil {
		return false, err
	}
	return isDataVolumeImportStatusSuccessful(dataVolume), nil

}
func IsDataSourceReady(dataSource *cdiv1beta1.DataSource) bool {
	return getConditionMapDs(dataSource)[cdiv1beta1.DataSourceReady].Status == v1.ConditionTrue
}

func HasDataVolumeFailedToImport(dataVolume *cdiv1beta1.DataVolume) bool {
	conditions := getConditionMapDv(dataVolume)
	return dataVolume.Status.Phase == cdiv1beta1.ImportInProgress &&
		dataVolume.Status.RestartCount > constants.UnusualRestartCountThreshold &&
		conditions[cdiv1beta1.DataVolumeRunning].Status == v1.ConditionFalse &&
		conditions[cdiv1beta1.DataVolumeRunning].Reason == constants.ReasonError
}

func isDataVolumeImportStatusSuccessful(dataVolume *cdiv1beta1.DataVolume) bool {
	return getConditionMapDv(dataVolume)[cdiv1beta1.DataVolumeBound].Status == v1.ConditionTrue &&
		dataVolume.Status.Phase == cdiv1beta1.Succeeded
}
