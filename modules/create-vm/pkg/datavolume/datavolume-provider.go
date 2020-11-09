package datavolume

import (
	"errors"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/k8s"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	datavolumev1alpha1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
	datavolumeclientv1alpha1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1alpha1"
)

type dataVolumeProvider struct {
	client datavolumeclientv1alpha1.CdiV1alpha1Interface
}

type DataVolumeProvider interface {
	GetByName(namespace string, names ...string) ([]*datavolumev1alpha1.DataVolume, error)
	AddOwnerReferences(dv *datavolumev1alpha1.DataVolume, newOwnerRefs ...metav1.OwnerReference) (*datavolumev1alpha1.DataVolume, error)
}

func NewDataVolumeProvider(client datavolumeclientv1alpha1.CdiV1alpha1Interface) DataVolumeProvider {
	return &dataVolumeProvider{
		client: client,
	}
}

func (d *dataVolumeProvider) GetByName(namespace string, names ...string) ([]*datavolumev1alpha1.DataVolume, error) {
	var multiError zerrors.MultiError
	var dvs []*datavolumev1alpha1.DataVolume

	for _, name := range names {
		dv, err := d.client.DataVolumes(namespace).Get(name, metav1.GetOptions{})
		if err == nil {
			dvs = append(dvs, dv)
		} else {
			dvs = append(dvs, nil)
			multiError.Add(name, err)
		}
	}
	return dvs, multiError.AsOptional()
}

func (d *dataVolumeProvider) AddOwnerReferences(dv *datavolumev1alpha1.DataVolume, newOwnerRefs ...metav1.OwnerReference) (*datavolumev1alpha1.DataVolume, error) {
	if dv == nil {
		return nil, errors.New("did not receive any DataVolume to add reference to")
	}

	if len(newOwnerRefs) <= 0 {
		return dv, nil
	}

	result := dv.DeepCopy()
	result.SetOwnerReferences(k8s.AppendOwnerReferences(result.GetOwnerReferences(), newOwnerRefs))

	patch, err := k8s.CreatePatch(dv, result)

	if err != nil {
		return nil, err
	}

	return d.client.DataVolumes(dv.Namespace).Patch(dv.Name, types.JSONPatchType, patch)
}
