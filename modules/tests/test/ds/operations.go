package ds

import (
	"context"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	cdiv1beta12 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	cdicliv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"
)

func WaitForSuccessfulDataSource(cdiClientSet cdicliv1beta1.CdiV1beta1Interface, namespace, name string, timeout time.Duration) error {
	return wait.PollImmediate(constants.PollInterval, timeout, func() (bool, error) {
		dataSource, err := cdiClientSet.DataSources(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return isDataSourceImportStatusSuccessful(dataSource), nil
	})
}

func IsDataSourceImportSuccessful(cdiClientSet cdicliv1beta1.CdiV1beta1Interface, namespace string, name string) bool {
	dataSource, err := cdiClientSet.DataSources(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return isDataSourceImportStatusSuccessful(dataSource)
}

func HasDataSourceFailedToImport(dataSource *cdiv1beta12.DataSource) bool {
	conditions := GetConditionMap(dataSource)
	return conditions[cdiv1beta12.DataSourceReady] == v1.ConditionFalse
}

func isDataSourceImportStatusSuccessful(dataSource *cdiv1beta12.DataSource) bool {
	return GetConditionMap(dataSource)[cdiv1beta12.DataSourceReady] == v1.ConditionTrue
}
