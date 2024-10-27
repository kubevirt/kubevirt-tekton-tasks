package secrets_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakek8sclient "k8s.io/client-go/kubernetes/fake"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/secrets"
)

var _ = Describe("VMExport", func() {
	const (
		namespace = "test-namespace"
		name      = "test-vmexport"
	)

	var (
		kubeClient *fakek8sclient.Clientset
	)

	BeforeEach(func() {
		kubeClient = fakek8sclient.NewSimpleClientset()

		os.Setenv("POD_NAME", name)
		os.Setenv("POD_NAMESPACE", namespace)
	})

	AfterEach(func() {
		os.Unsetenv("POD_NAME")
		os.Unsetenv("POD_NAMESPACE")
	})

	Describe("CreateVirtualMachineExportSecret", func() {
		It("should return no error when created", func() {
			_, err := kubeClient.CoreV1().Pods(namespace).Create(context.Background(),
				&corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
				},
				metav1.CreateOptions{},
			)
			Expect(err).NotTo(HaveOccurred())

			err = secrets.CreateVirtualMachineExportSecret(kubeClient, namespace, name)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when set pod owner reference failed", func() {
			err := secrets.CreateVirtualMachineExportSecret(kubeClient, namespace, name)
			Expect(err).To(MatchError(errors.IsNotFound, "errors.IsNotFound"))
		})
	})

	Describe("GetTokenFromVirtualMachineExportSecret", func() {
		It("should return token when exists", func() {
			_, err := kubeClient.CoreV1().Secrets(namespace).Create(context.Background(),
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Data: map[string][]byte{
						"token": []byte("fake"),
					},
				},
				metav1.CreateOptions{},
			)
			Expect(err).NotTo(HaveOccurred())

			token, err := secrets.GetTokenFromVirtualMachineExportSecret(kubeClient, namespace, name)
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal("fake"))
		})

		It("should return error when not found", func() {
			_, err := secrets.GetTokenFromVirtualMachineExportSecret(kubeClient, namespace, name)
			Expect(err).To(MatchError(errors.IsNotFound, "errors.IsNotFound"))
		})

		It("should return error when failed to get token", func() {
			_, err := kubeClient.CoreV1().Secrets(namespace).Create(context.Background(),
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
				},
				metav1.CreateOptions{},
			)
			Expect(err).NotTo(HaveOccurred())

			_, err = secrets.GetTokenFromVirtualMachineExportSecret(kubeClient, namespace, name)
			Expect(err).To(MatchError("failed to get export token from 'test-namespace/test-vmexport'"))
		})
	})
})
