package k8s

import "k8s.io/apimachinery/pkg/apis/meta/v1"

func AppendOwnerReferences(ownerRefs []v1.OwnerReference, newOwnerRefs []v1.OwnerReference) []v1.OwnerReference {
	if ownerRefs == nil {
		ownerRefs = []v1.OwnerReference{}
	}

	for _, newOwnerRef := range newOwnerRefs {
		ownerRefs = append(ownerRefs, newOwnerRef)
	}
	return ownerRefs
}
