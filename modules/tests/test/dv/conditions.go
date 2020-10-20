package dv

import (
	v1 "k8s.io/api/core/v1"
	"kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
)

func GetConditionMap(dv *v1beta1.DataVolume) map[v1beta1.DataVolumeConditionType]v1.ConditionStatus {
	result := make(map[v1beta1.DataVolumeConditionType]v1.ConditionStatus)
	for _, cond := range dv.Status.Conditions {
		result[cond.Type] = cond.Status
	}
	return result
}
