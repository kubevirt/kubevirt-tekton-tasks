package vmexport_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/golang/mock/gomock"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakek8sclient "k8s.io/client-go/kubernetes/fake"

	v1beta1 "kubevirt.io/api/export/v1beta1"
	fakecdiclient "kubevirt.io/client-go/generated/containerized-data-importer/clientset/versioned/fake"
	kubevirtfake "kubevirt.io/client-go/generated/kubevirt/clientset/versioned/fake"
	"kubevirt.io/client-go/kubecli"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/vmexport"
)

var _ = Describe("VMExport", func() {
	const (
		namespace = "test-namespace"
		name      = "test-vmexport"
	)

	var (
		kubeClient     *fakek8sclient.Clientset
		cdiClient      *fakecdiclient.Clientset
		vmExportClient *kubevirtfake.Clientset
		virtClient     kubecli.KubevirtClient
	)

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		kubeClient = fakek8sclient.NewSimpleClientset()
		vmExportClient = kubevirtfake.NewSimpleClientset()
		cdiClient = fakecdiclient.NewSimpleClientset()

		kubecli.GetKubevirtClientFromClientConfig = kubecli.GetMockKubevirtClientFromClientConfig
		kubecli.MockKubevirtClientInstance = kubecli.NewMockKubevirtClient(ctrl)
		kubecli.MockKubevirtClientInstance.EXPECT().CoreV1().Return(kubeClient.CoreV1()).AnyTimes()
		kubecli.MockKubevirtClientInstance.EXPECT().VirtualMachineExport(namespace).Return(vmExportClient.ExportV1beta1().VirtualMachineExports(namespace)).AnyTimes()
		kubecli.MockKubevirtClientInstance.EXPECT().CdiClient().Return(cdiClient).AnyTimes()

		virtClient, _ = kubecli.GetKubevirtClientFromClientConfig(nil)

		os.Setenv("POD_NAME", name)
		os.Setenv("POD_NAMESPACE", namespace)
	})

	AfterEach(func() {
		os.Unsetenv("POD_NAME")
		os.Unsetenv("POD_NAMESPACE")
	})

	Describe("CreateVirtualMachineExport", func() {
		DescribeTable("should return no error when created",
			func(resource string) {
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

				err = vmexport.CreateVirtualMachineExport(virtClient, resource, namespace, name)
				Expect(err).NotTo(HaveOccurred())
				Expect(errors.IsNotFound(err)).To(BeFalse())
			},
			Entry("export-source-kind is vm", "vm"),
			Entry("export-source-kind is vmsnapshot", "vmsnapshot"),
			Entry("export-source-kind is pvc", "pvc"),
		)

		It("should return error when export-source-kind invalid", func() {
			err := vmexport.CreateVirtualMachineExport(virtClient, "fake", namespace, name)
			Expect(err).To(MatchError("invalid export-source-kind: fake, must be one of vm, vmsnapshot, pvc"))
		})

		It("should return error when set pod owner reference failed", func() {
			err := vmexport.CreateVirtualMachineExport(virtClient, "vm", namespace, name)
			Expect(err).To(MatchError(errors.IsNotFound, "errors.IsNotFound"))
		})
	})

	Describe("WaitUntilVirtualMachineExportReady", func() {
		It("should return no error", func() {
			_, err := vmExportClient.ExportV1beta1().VirtualMachineExports(namespace).Create(context.Background(),
				&v1beta1.VirtualMachineExport{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Status: &v1beta1.VirtualMachineExportStatus{
						Phase: v1beta1.Ready,
					},
				},
				metav1.CreateOptions{},
			)
			Expect(err).NotTo(HaveOccurred())

			err = vmexport.WaitUntilVirtualMachineExportReady(virtClient, namespace, name)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetRawDiskUrlFromVolumes", func() {
		Context("when retrieved URL from the VirtualMachineExport volumes", func() {
			BeforeEach(func() {
				_, err := vmExportClient.ExportV1beta1().VirtualMachineExports(namespace).Create(context.Background(),
					&v1beta1.VirtualMachineExport{
						ObjectMeta: metav1.ObjectMeta{
							Name:      name,
							Namespace: namespace,
						},
						Status: &v1beta1.VirtualMachineExportStatus{
							Links: &v1beta1.VirtualMachineExportLinks{
								Internal: &v1beta1.VirtualMachineExportLink{
									Volumes: []v1beta1.VirtualMachineExportVolume{
										{
											Name: "disk-volume-0",
											Formats: []v1beta1.VirtualMachineExportVolumeFormat{
												{
													Format: v1beta1.KubeVirtRaw,
													Url:    "https://vmexport-proxy.test.net/volumes/disk-volume-0/disk.img",
												},
											},
										},
										{
											Name: "disk-volume-1",
											Formats: []v1beta1.VirtualMachineExportVolumeFormat{
												{
													Format: v1beta1.KubeVirtRaw,
													Url:    "https://vmexport-proxy.test.net/volumes/disk-volume-1/disk.img",
												},
											},
										},
									},
								},
							},
						},
					},
					metav1.CreateOptions{},
				)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return URL correctly", func() {
				url, err := vmexport.GetRawDiskUrlFromVolumes(virtClient, namespace, name, "disk-volume-1")
				Expect(err).NotTo(HaveOccurred())
				Expect(url).To(Equal("https://vmexport-proxy.test.net/volumes/disk-volume-1/disk.img"))
			})

			It("should return error when no volume found", func() {
				_, err := vmexport.GetRawDiskUrlFromVolumes(virtClient, namespace, name, "disk-volume-2")
				Expect(err).To(MatchError("volume disk-volume-2 is not found in VirtualMachineExport internal volumes"))
			})
		})

		It("should return not found error", func() {
			_, err := vmexport.GetRawDiskUrlFromVolumes(virtClient, namespace, name, "disk-volume-1")
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

			_, err = vmexport.GetRawDiskUrlFromVolumes(virtClient, namespace, name, "disk-volume-1")
			Expect(err).To(MatchError("no links found in VirtualMachineExport status"))
		})
	})

	Describe("GetLabelsFromExportSource", func() {
		var (
			defaultInstanceType = "instancetype.kubevirt.io/default-instancetype"
			defaultPreference   = "instancetype.kubevirt.io/default-preference"
		)

		DescribeTable("should return labels from datavolume", func(resourceType string) {
			_, err := virtClient.CdiClient().CdiV1beta1().DataVolumes(namespace).Create(context.Background(), &cdiv1beta1.DataVolume{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-volume",
					Namespace: namespace,
					Labels: map[string]string{
						defaultInstanceType: "cx1.2xlarge",
						defaultPreference:   "fedora",
					},
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			labels, err := vmexport.GetLabelsFromExportSource(virtClient, resourceType, namespace, name, "example-volume")
			Expect(err).NotTo(HaveOccurred())
			Expect(labels).To(HaveKeyWithValue(defaultInstanceType, "cx1.2xlarge"))
			Expect(labels).To(HaveKeyWithValue(defaultPreference, "fedora"))
		},
			Entry("for VirtualMachine", "VirtualMachine"),
			Entry("for VirtualMachineSnapshot", "VirtualMachineSnapshot"),
		)

		DescribeTable("should return labels from pvc if there is no datavolume", func(resourceType string) {
			_, err := kubeClient.CoreV1().PersistentVolumeClaims(namespace).Create(context.Background(), &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-pvc",
					Namespace: namespace,
					Labels: map[string]string{
						defaultInstanceType: "cx1.2xlarge",
						defaultPreference:   "fedora",
					},
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			labels, err := vmexport.GetLabelsFromExportSource(virtClient, resourceType, namespace, name, "example-pvc")
			Expect(err).NotTo(HaveOccurred())
			Expect(labels).To(HaveKeyWithValue(defaultInstanceType, "cx1.2xlarge"))
			Expect(labels).To(HaveKeyWithValue(defaultPreference, "fedora"))
		},
			Entry("for VirtualMachine", "VirtualMachine"),
			Entry("for VirtualMachineSnapshot", "VirtualMachineSnapshot"),
		)

		It("should return labels directly from pvc", func() {
			_, err := kubeClient.CoreV1().PersistentVolumeClaims(namespace).Create(context.Background(), &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-pvc",
					Namespace: namespace,
					Labels: map[string]string{
						defaultInstanceType: "cx1.2xlarge",
						defaultPreference:   "fedora",
					},
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			labels, err := vmexport.GetLabelsFromExportSource(virtClient, "PersistentVolumeClaim", namespace, "example-pvc", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(labels).To(HaveKeyWithValue(defaultInstanceType, "cx1.2xlarge"))
			Expect(labels).To(HaveKeyWithValue(defaultPreference, "fedora"))
		})

		It("should return unsupported source kind error", func() {
			labels, err := vmexport.GetLabelsFromExportSource(virtClient, "VirtualMachineInstance", namespace, name, "disk-volume-1")
			Expect(err).To(MatchError("unsupported source kind: VirtualMachineInstance"))
			Expect(labels).To(BeNil())
		})

		DescribeTable("should return an error when DV or PVC is not found", func(exportSourceKind string, volumeName string) {
			labels, err := vmexport.GetLabelsFromExportSource(virtClient, exportSourceKind, namespace, name, volumeName)
			Expect(err).To(MatchError(errors.IsNotFound, "errors.IsNotFound"))
			Expect(labels).To(BeNil())
		},
			Entry("for VirtualMachine with missing DV and PVC", "VirtualMachine", "example-volume"),
			Entry("for VirtualMachineSnapshot with missing DV and PVC", "VirtualMachineSnapshot", "example-volume"),
			Entry("for PersistentVolumeClaim with missing PVC", "PersistentVolumeClaim", ""),
		)
	})
})
