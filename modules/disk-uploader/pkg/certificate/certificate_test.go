package certificate_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/golang/mock/gomock"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1beta1 "kubevirt.io/api/export/v1beta1"
	kubevirtfake "kubevirt.io/client-go/generated/kubevirt/clientset/versioned/fake"
	"kubevirt.io/client-go/kubecli"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/certificate"
)

var _ = Describe("Certificate", func() {
	const (
		namespace = "test-namespace"
		name      = "test-vmexport"
	)

	var (
		vmExportClient *kubevirtfake.Clientset
		virtClient     kubecli.KubevirtClient
	)

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		vmExportClient = kubevirtfake.NewSimpleClientset()

		kubecli.GetKubevirtClientFromClientConfig = kubecli.GetMockKubevirtClientFromClientConfig
		kubecli.MockKubevirtClientInstance = kubecli.NewMockKubevirtClient(ctrl)
		kubecli.MockKubevirtClientInstance.EXPECT().VirtualMachineExport(namespace).Return(vmExportClient.ExportV1beta1().VirtualMachineExports(namespace)).AnyTimes()

		virtClient, _ = kubecli.GetKubevirtClientFromClientConfig(nil)
	})

	Describe("GetCertificateFromVirtualMachineExport", func() {
		It("should return the certificate content", func() {
			_, err := vmExportClient.ExportV1beta1().VirtualMachineExports(namespace).Create(context.Background(),
				&v1beta1.VirtualMachineExport{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Status: &v1beta1.VirtualMachineExportStatus{
						Links: &v1beta1.VirtualMachineExportLinks{
							Internal: &v1beta1.VirtualMachineExportLink{
								Cert: "test-cert-content",
							},
						},
					},
				},
				metav1.CreateOptions{},
			)
			Expect(err).NotTo(HaveOccurred())

			cert, err := certificate.GetCertificateFromVirtualMachineExport(virtClient, namespace, name)
			Expect(err).NotTo(HaveOccurred())
			Expect(cert).To(Equal("test-cert-content"))
		})

		It("should return not found error", func() {
			_, err := certificate.GetCertificateFromVirtualMachineExport(virtClient, namespace, "test")
			Expect(err).To(MatchError(errors.IsNotFound, "errors.IsNotFound"))
		})

		It("should return an error when no links found", func() {
			_, err := vmExportClient.ExportV1beta1().VirtualMachineExports(namespace).Create(context.Background(),
				&v1beta1.VirtualMachineExport{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Status: &v1beta1.VirtualMachineExportStatus{
						Links: nil,
					},
				},
				metav1.CreateOptions{},
			)
			Expect(err).NotTo(HaveOccurred())

			cert, err := certificate.GetCertificateFromVirtualMachineExport(virtClient, namespace, name)
			Expect(err).To(MatchError("no links found in VirtualMachineExport status"))
			Expect(cert).To(BeEmpty())
		})

		It("should return an error when no certificate found", func() {
			_, err := vmExportClient.ExportV1beta1().VirtualMachineExports(namespace).Create(context.Background(),
				&v1beta1.VirtualMachineExport{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Status: &v1beta1.VirtualMachineExportStatus{
						Links: &v1beta1.VirtualMachineExportLinks{
							Internal: &v1beta1.VirtualMachineExportLink{
								Cert: "",
							},
						},
					},
				},
				metav1.CreateOptions{},
			)
			Expect(err).NotTo(HaveOccurred())

			cert, err := certificate.GetCertificateFromVirtualMachineExport(virtClient, namespace, name)
			Expect(err).To(MatchError("no certificate found in VirtualMachineExport status"))
			Expect(cert).To(BeEmpty())
		})
	})
})
