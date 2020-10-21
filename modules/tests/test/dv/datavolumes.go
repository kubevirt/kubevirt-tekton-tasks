package dv

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta12 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
)

type TestDataVolume struct {
	Data v1beta12.DataVolume
}

func NewBlankDataVolume(name string) *TestDataVolume {
	volumeMode := v1.PersistentVolumeFilesystem
	datavolume := v1beta12.DataVolume{
		TypeMeta: metav1.TypeMeta{
			APIVersion: constants.DataVolumeApiVersion,
			Kind:       constants.DataVolumeKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1beta12.DataVolumeSpec{
			PVC: &v1.PersistentVolumeClaimSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceStorage: *resource.NewScaledQuantity(100, resource.Mega),
					},
				},
				VolumeMode: &volumeMode,
			},
			Source: v1beta12.DataVolumeSource{
				Blank: &v1beta12.DataVolumeBlankImage{},
			},
		},
	}

	return &TestDataVolume{
		datavolume,
	}
}

func (d *TestDataVolume) WithoutTypeMeta() *TestDataVolume {
	d.Data.Kind = ""
	d.Data.APIVersion = ""
	return d
}
