package secrets

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/ownerreference"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateVirtualMachineExportSecret(k8sClient kubernetes.Interface, namespace, baseName string) (*corev1.Secret, error) {
	length := 20
	token, err := generateSecureRandomString(length)
	if err != nil {
		return nil, err
	}

	v1Secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: baseName + "-",
			Namespace:    namespace,
		},
		StringData: map[string]string{
			"token": token,
		},
	}

	if err := ownerreference.SetPodOwnerReference(k8sClient, v1Secret); err != nil {
		return nil, err
	}

	return k8sClient.CoreV1().Secrets(namespace).Create(context.Background(), v1Secret, metav1.CreateOptions{})
}

func GetTokenFromVirtualMachineExportSecret(k8sClient kubernetes.Interface, namespace, name string) (string, error) {
	secret, err := k8sClient.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	data := secret.Data["token"]
	if len(data) == 0 {
		return "", fmt.Errorf("failed to get export token from '%s/%s'", namespace, name)
	}
	return string(data), nil
}

func generateSecureRandomString(n int) (string, error) {
	// Alphanums is the list of alphanumeric characters used to create a securely generated random string
	alphanums := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	ret := make([]byte, n)
	for i := range ret {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphanums))))
		if err != nil {
			return "", err
		}
		ret[i] = alphanums[num.Int64()]
	}

	return string(ret), nil
}
