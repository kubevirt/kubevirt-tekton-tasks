package ownerreference

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	podNameEnv      = "POD_NAME"
	podNamespaceEnv = "POD_NAMESPACE"
)

func getTaskRunPod(k8sClient *k8sv1.CoreV1Client) (*corev1.Pod, error) {
	podName, isSet := os.LookupEnv(podNameEnv)
	if !isSet {
		return nil, fmt.Errorf("pod name env variable is not set")
	}

	podNamespace, isSet := os.LookupEnv(podNamespaceEnv)
	if !isSet {
		return nil, fmt.Errorf("pod namespace env variable is not set")
	}

	pod := &corev1.Pod{}
	pod, err := k8sClient.Pods(podNamespace).Get(context.Background(), podName, metav1.GetOptions{})
	return pod, err
}

func SetPodOwnerReference(k8sClient *k8sv1.CoreV1Client, object metav1.Object) error {
	pod, err := getTaskRunPod(k8sClient)
	if err != nil {
		return err
	}

	if object.GetNamespace() != pod.GetNamespace() {
		return fmt.Errorf("can't create owner reference for objects in different namespaces")
	}

	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)

	gvks, _, err := scheme.ObjectKinds(pod)
	if err != nil {
		return fmt.Errorf("could not get GroupVersionKind for object: %w", err)
	}
	ref := metav1.OwnerReference{
		APIVersion: gvks[0].GroupVersion().String(),
		Kind:       gvks[0].Kind,
		UID:        pod.GetUID(),
		Name:       pod.GetName(),
	}

	object.SetOwnerReferences([]metav1.OwnerReference{ref})
	return nil

}
