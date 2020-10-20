package dv

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/tests/test/constants"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta12 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
)

type DV struct {
	datavolume *v1beta12.DataVolume
}

func NewBlankDV(name string) *DV {
	volumeMode := v1.PersistentVolumeFilesystem
	datavolume := &v1beta12.DataVolume{
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

	return &DV{
		datavolume,
	}
}

func (d *DV) WithoutTypeMeta() *DV {
	d.datavolume.Kind = ""
	d.datavolume.APIVersion = ""
	return d
}

func (d *DV) Build() *v1beta12.DataVolume {
	return d.datavolume
}
