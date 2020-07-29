package pvc

import (
	"errors"
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/k8s"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type pvcProvider struct {
	client clientv1.CoreV1Interface
}

type PersistentVolumeClaimProvider interface {
	GetByName(namespace string, names ...string) ([]*v1.PersistentVolumeClaim, error)
	AddOwnerReferences(dv *v1.PersistentVolumeClaim, newOwnerRefs ...metav1.OwnerReference) (*v1.PersistentVolumeClaim, error)
}

func NewPersistentVolumeClaimProvider(client clientv1.CoreV1Interface) PersistentVolumeClaimProvider {
	return &pvcProvider{
		client: client,
	}
}

func (d *pvcProvider) GetByName(namespace string, names ...string) ([]*v1.PersistentVolumeClaim, error) {
	var errs []error
	var pvcs []*v1.PersistentVolumeClaim

	for _, name := range names {
		pvc, err := d.client.PersistentVolumeClaims(namespace).Get(name, metav1.GetOptions{})
		if err == nil {
			pvcs = append(pvcs, pvc)
			errs = append(errs, nil)
		} else {
			pvcs = append(pvcs, nil)
			errs = append(errs, err)
		}
	}
	return pvcs, errors2.NewMultiErrorOrNil(errs)
}

func (d *pvcProvider) AddOwnerReferences(pvc *v1.PersistentVolumeClaim, newOwnerRefs ...metav1.OwnerReference) (*v1.PersistentVolumeClaim, error) {
	if pvc == nil {
		return nil, errors.New("did not receive any PersistentVolumeClaim to add reference to")
	}

	if len(newOwnerRefs) <= 0 {
		return pvc, nil
	}

	result := pvc.DeepCopy()
	result.SetOwnerReferences(k8s.AppendOwnerReferences(result.GetOwnerReferences(), newOwnerRefs))

	patch, err := k8s.CreatePatch(pvc, result)

	if err != nil {
		return nil, err
	}

	return d.client.PersistentVolumeClaims(pvc.Namespace).Patch(pvc.Name, types.JSONPatchType, patch)
}
