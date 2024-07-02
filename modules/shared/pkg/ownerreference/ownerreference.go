package ownerreference

import (
	"context"
	"fmt"
	"os"

	apiv1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	podNameEnv      = "POD_NAME"
	podNamespaceEnv = "POD_NAMESPACE"
)

func getTaskRunPod(k8sClient client.Client) (*apiv1.Pod, error) {
	podName, isSet := os.LookupEnv(podNameEnv)
	if !isSet {
		return nil, fmt.Errorf("pod name env variable is not set")
	}

	podNamespace, isSet := os.LookupEnv(podNamespaceEnv)
	if !isSet {
		return nil, fmt.Errorf("pod namespace env variable is not set")
	}
	objKey := client.ObjectKey{Namespace: podNamespace, Name: podName}
	pod := &apiv1.Pod{}
	err := k8sClient.Get(context.Background(), objKey, pod)
	return pod, err
}

func SetPodOwnerReference(k8sClient client.Client, object client.Object) error {
	pod, err := getTaskRunPod(k8sClient)
	if err != nil {
		return err
	}

	if object.GetNamespace() != pod.GetNamespace() {
		return fmt.Errorf("can't create owner reference for objects in different namespaces")
	}

	return controllerutil.SetOwnerReference(pod, object, k8sClient.Scheme())

}
