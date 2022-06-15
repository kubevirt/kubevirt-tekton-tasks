package ds

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta12 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

const (
	dataSourceKind       = "DataSource"
	dataSourceApiVersion = "cdi.kubevirt.io/v1beta1"
)

func NewDataSource(name string) *v1beta12.DataSource {
	dataSource := &v1beta12.DataSource{
		TypeMeta: metav1.TypeMeta{
			APIVersion: dataSourceApiVersion,
			Kind:       dataSourceKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1beta12.DataSourceSpec{
			Source: v1beta12.DataSourceSource{
				PVC: &v1beta12.DataVolumeSourcePVC{
					Name: name,
				},
			},
		},
	}

	return dataSource
}
