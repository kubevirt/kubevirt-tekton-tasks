package ds

import (
	v1 "k8s.io/api/core/v1"
	"kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func GetConditionMap(ds *v1beta1.DataSource) map[v1beta1.DataSourceConditionType]v1.ConditionStatus {
	result := make(map[v1beta1.DataSourceConditionType]v1.ConditionStatus)
	for _, cond := range ds.Status.Conditions {
		result[cond.Type] = cond.Status
	}
	return result
}
