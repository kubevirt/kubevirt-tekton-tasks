package datasource

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/yaml"
)

const (
	dataSourceKind       = "DataSource"
	dataSourceApiVersion = "cdi.kubevirt.io/v1beta1"
)

type TestDataSource struct {
	Data *cdiv1beta1.DataSource
}

func NewDataSource(name string) *TestDataSource {
	dataSource := &cdiv1beta1.DataSource{
		TypeMeta: metav1.TypeMeta{
			APIVersion: dataSourceApiVersion,
			Kind:       dataSourceKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: cdiv1beta1.DataSourceSpec{
			Source: cdiv1beta1.DataSourceSource{
				PVC: &cdiv1beta1.DataVolumeSourcePVC{},
			},
		},
	}

	return &TestDataSource{
		dataSource,
	}
}

func (d *TestDataSource) WithoutTypeMeta() *TestDataSource {
	d.Data.Kind = ""
	d.Data.APIVersion = ""
	return d
}

func (d *TestDataSource) WithAPIVersion(apiVersion string) *TestDataSource {
	d.Data.APIVersion = apiVersion
	return d
}

func (d *TestDataSource) WithKind(kind string) *TestDataSource {
	d.Data.Kind = kind
	return d
}

func (d *TestDataSource) WithNamespace(namespace string) *TestDataSource {
	d.Data.Namespace = namespace
	return d
}

func (d *TestDataSource) WithSourcePVC(name, namespace string) *TestDataSource {
	d.Data.Spec.Source.PVC = &cdiv1beta1.DataVolumeSourcePVC{
		Name:      name,
		Namespace: namespace,
	}
	return d
}

func (d *TestDataSource) WithGenerateName(generateName string) *TestDataSource {
	d.Data.GenerateName = generateName
	return d
}

func (d *TestDataSource) Build() *cdiv1beta1.DataSource {
	return d.Data
}

func (t *TestDataSource) ToString() string {
	outBytes, _ := yaml.Marshal(t.Data)
	return string(outBytes)
}
