package datavolume

import (
	"context"
	"errors"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/k8s"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	datavolumev1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	datavolumeclientv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"
)

type dataVolumeProvider struct {
	client datavolumeclientv1beta1.CdiV1beta1Interface
}

type DataVolumeProvider interface {
	GetByName(namespace string, names ...string) ([]*datavolumev1beta1.DataVolume, error)
	AddOwnerReferences(dv *datavolumev1beta1.DataVolume, newOwnerRefs ...metav1.OwnerReference) (*datavolumev1beta1.DataVolume, error)
}

func NewDataVolumeProvider(client datavolumeclientv1beta1.CdiV1beta1Interface) DataVolumeProvider {
	return &dataVolumeProvider{
		client: client,
	}
}

func (d *dataVolumeProvider) GetByName(namespace string, names ...string) ([]*datavolumev1beta1.DataVolume, error) {
	var multiError zerrors.MultiError
	var dvs []*datavolumev1beta1.DataVolume

	for _, name := range names {
		dv, err := d.client.DataVolumes(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err == nil {
			dvs = append(dvs, dv)
		} else {
			dvs = append(dvs, nil)
			multiError.Add(name, err)
		}
	}
	return dvs, multiError.AsOptional()
}

func (d *dataVolumeProvider) AddOwnerReferences(dv *datavolumev1beta1.DataVolume, newOwnerRefs ...metav1.OwnerReference) (*datavolumev1beta1.DataVolume, error) {
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

	return d.client.DataVolumes(dv.Namespace).Patch(context.TODO(), dv.Name, types.JSONPatchType, patch, metav1.PatchOptions{})
}
