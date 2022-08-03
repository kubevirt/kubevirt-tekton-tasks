package dataobject

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-data-object/pkg/constants"
	v1 "k8s.io/api/core/v1"
	"kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func hasDataVolumeFailedToImport(dv *v1beta1.DataVolume) bool {
	conditions := getConditionMapDv(dv)
	return dv.Status.Phase == v1beta1.ImportInProgress &&
		dv.Status.RestartCount > constants.UnusualRestartCountThreshold &&
		conditions[v1beta1.DataVolumeBound].Status == v1.ConditionTrue &&
		conditions[v1beta1.DataVolumeRunning].Status == v1.ConditionFalse &&
		conditions[v1beta1.DataVolumeRunning].Reason == constants.ReasonError
}

func isDataVolumeImportStatusSuccessful(dv *v1beta1.DataVolume) bool {
	conditions := getConditionMapDv(dv)
	return dv.Status.Phase == v1beta1.Succeeded &&
		conditions[v1beta1.DataVolumeBound].Status == v1.ConditionTrue
}

func isDataSourceReady(dataSource *v1beta1.DataSource) bool {
	return getConditionMapDs(dataSource)[v1beta1.DataSourceReady].Status == v1.ConditionTrue
}
