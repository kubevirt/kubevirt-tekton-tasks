package k8s

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func EnsureLabels(objectMeta *v1.ObjectMeta) map[string]string {
	labels := objectMeta.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
		objectMeta.SetLabels(labels)
	}
	return objectMeta.Labels
}

func EnsureAnnotations(objectMeta *v1.ObjectMeta) map[string]string {
	annotations := objectMeta.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
		objectMeta.SetAnnotations(annotations)
	}
	return objectMeta.Annotations
}
