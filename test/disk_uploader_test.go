package test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/ptr"

	v1 "kubevirt.io/api/core/v1"
	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"kubevirt.io/kubevirt/pkg/libvmi"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/framework"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/runner"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/testconfigs"
	"github.com/kubevirt/kubevirt-tekton-tasks/test/vm"
)

var (
	secretName       string
	imageDestination string
	registryKeyId    string
	registryKey      string
)

var _ = Describe("Run disk-uploader", func() {
	f := framework.NewFramework()

	BeforeEach(func() {
		if f.TestOptions.SkipDiskUploaderTests {
			Skip("skipDiskUploaderTests is set to true, skipping tests")
		}

		imageDestination = os.Getenv("IMAGE_DESTINATION")
		Expect(imageDestination).ToNot(BeEmpty())

		registryKeyId = os.Getenv("REGISTRY_ACCESS_KEY_ID")
		Expect(registryKeyId).ToNot(BeEmpty())

		registryKey = os.Getenv("REGISTRY_SECRET_KEY")
		Expect(registryKey).ToNot(BeEmpty())

		secretName = constants.E2ETestsRandomName("disk-uploader-credentials")
		_, err := createRegistryCredentials(f.CoreV1Client, secretName, f.DeployNamespace)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Extracts disk from VM and upload to container registry", func() {
		alpineDataVolume := newAlpineDataVolume()
		alpineVm, err := createAlpineVM(f.KubevirtClient, f.DeployNamespace, alpineDataVolume)
		Expect(err).ToNot(HaveOccurred())

		f.ManageVMs(alpineVm)
		f.ManageDataVolumes(alpineDataVolume)

		_, err = vm.WaitForVM(f.KubevirtClient, f.DeployNamespace, alpineVm.Name, "", constants.Timeouts.WaitForVMStart.Duration, false)
		Expect(err).ToNot(HaveOccurred())

		config := &testconfigs.DiskUploaderTestConfig{
			TaskRunTestConfig: testconfigs.TaskRunTestConfig{},
			TaskData: testconfigs.DiskUploaderTaskData{
				ExportSourceKind: "vm",
				ExportSourceName: alpineVm.Name,
				VolumeName:       alpineVm.Spec.Template.Spec.Volumes[0].DataVolume.Name,
				ImageDestination: imageDestination,
				SecretName:       secretName,
			},
		}
		f.TestSetup(config)

		taskRun := runner.NewTaskRunRunner(f, config.GetTaskRun()).
			CreateTaskRun().
			ExpectSuccess().
			ExpectLogs(config.GetAllExpectedLogs()...)

		digest := taskRun.GetResults()[constants.DigestResultName]
		Expect(digest).ToNot(BeEmpty())

		ref, err := name.ParseReference(imageDestination)
		Expect(err).ToNot(HaveOccurred())

		descriptor, err := remote.Get(ref, remote.WithAuth(&authn.Basic{
			Username: registryKeyId,
			Password: registryKey,
		}))
		Expect(err).ToNot(HaveOccurred())
		Expect(descriptor.Digest.String()).To(Equal(digest))
	})
})

func createRegistryCredentials(client kubeclientcorev1.CoreV1Interface, name, namespace string) (*corev1.Secret, error) {
	v1Secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"accessKeyId": []byte(registryKeyId),
			"secretKey":   []byte(registryKey),
		},
	}

	return client.Secrets(namespace).Create(context.Background(), v1Secret, metav1.CreateOptions{})
}

func createAlpineVM(client kubevirtcliv1.KubevirtClient, namespace string, dataVolume *cdiv1beta1.DataVolume) (*v1.VirtualMachine, error) {
	v1VirtualMachine := libvmi.NewVirtualMachine(
		libvmi.New(
			libvmi.WithDataVolume("datavolumedisk", dataVolume.Name),
			libvmi.WithResourceMemory("256M"),
		),
		libvmi.WithDataVolumeTemplate(dataVolume),
	)

	return client.VirtualMachine(namespace).Create(context.Background(), v1VirtualMachine, metav1.CreateOptions{})
}

func newAlpineDataVolume() *cdiv1beta1.DataVolume {
	return &cdiv1beta1.DataVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: constants.E2ETestsRandomName("test-datavolume"),
			Annotations: map[string]string{
				"cdi.kubevirt.io/storage.bind.immediate.requested": "true",
			},
		},
		Spec: cdiv1beta1.DataVolumeSpec{
			Source: &cdiv1beta1.DataVolumeSource{
				Registry: &cdiv1beta1.DataVolumeSourceRegistry{
					URL: ptr.To("docker://quay.io/kubevirt/alpine-container-disk-demo:20250818_82ae6622ba"),
				},
			},
			Storage: &cdiv1beta1.StorageSpec{
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						"storage": resource.MustParse("512Mi"),
					},
				},
			},
		},
	}
}
