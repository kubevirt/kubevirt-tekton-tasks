package testobjects

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testconstants"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewTestSecret(stringData map[string]string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      "testsecret",
			Namespace: testconstants.NamespaceTestDefault,
		},
		StringData: stringData,
	}
}
