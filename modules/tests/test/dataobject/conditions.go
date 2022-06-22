package dataobject

import (
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func getConditionMapDv(dv *cdiv1beta1.DataVolume) map[cdiv1beta1.DataVolumeConditionType]cdiv1beta1.DataVolumeCondition {
	result := map[cdiv1beta1.DataVolumeConditionType]cdiv1beta1.DataVolumeCondition{}
	for _, cond := range dv.Status.Conditions {
		result[cond.Type] = cond
	}
	return result
}

func getConditionMapDs(ds *cdiv1beta1.DataSource) map[cdiv1beta1.DataSourceConditionType]cdiv1beta1.DataSourceCondition {
	result := map[cdiv1beta1.DataSourceConditionType]cdiv1beta1.DataSourceCondition{}
	for _, cond := range ds.Status.Conditions {
		result[cond.Type] = cond
	}
	return result
}
