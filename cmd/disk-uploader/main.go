package main

import (
	"os"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/certificate"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/disk"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/image"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/secrets"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/vmexport"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	res "github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/results"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	goarg "github.com/alexflint/go-arg"
	"go.uber.org/zap"
	kubecli "kubevirt.io/client-go/kubecli"
)

const (
	diskPath            string = "./tmp/disk.qcow2"
	certificatePath     string = "./tmp/tls.crt"
	kvExportTokenHeader string = "x-kubevirt-export-token"
)

func run(opts parse.CLIOptions, k8sClient kubernetes.Interface, virtClient kubecli.KubevirtClient) (string, error) {
	kind := opts.GetExportSourceKind()
	name := opts.GetExportSourceName()
	namespace := opts.GetExportSourceNamespace()
	volumeName := opts.GetVolumeName()
	imageDestination := opts.GetImageDestination()
	imagePushTimeout := opts.GetPushTimeout()

	log.Logger().Info("Creating a new Secret object...", zap.String("namespace", namespace), zap.String("name", name))

	vmExportSecret, err := secrets.CreateVirtualMachineExportSecret(k8sClient, namespace, name)
	if err != nil {
		return "", err
	}

	log.Logger().Info("Creating a new VirtualMachineExport object...", zap.String("namespace", namespace), zap.String("name", name))

	vmExport, err := vmexport.CreateVirtualMachineExport(virtClient, kind, namespace, name, vmExportSecret.Name)
	if err != nil {
		return "", err
	}

	log.Logger().Info("Waiting for VirtualMachineExport object status...")

	if err := vmexport.WaitUntilVirtualMachineExportReady(virtClient, namespace, vmExport.Name); err != nil {
		return "", err
	}

	log.Logger().Info("Getting raw disk URL from the VirtualMachineExport object status...")

	rawDiskUrl, err := vmexport.GetRawDiskUrlFromVolumes(virtClient, namespace, vmExport.Name, volumeName)
	if err != nil {
		return "", err
	}

	log.Logger().Info("Creating TLS certificate file from the VirtualMachineExport object status...")

	certificateData, err := certificate.GetCertificateFromVirtualMachineExport(virtClient, namespace, vmExport.Name)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(certificatePath, []byte(certificateData), 0644); err != nil {
		return "", err
	}

	log.Logger().Info("Getting export token from the Secret object...")

	kvExportToken, err := secrets.GetTokenFromVirtualMachineExportSecret(virtClient, namespace, vmExportSecret.Name)
	if err != nil {
		return "", err
	}

	log.Logger().Info("Downloading disk image from the VirtualMachineExport server...")

	if err := disk.DownloadDiskImageFromURL(rawDiskUrl, kvExportTokenHeader, kvExportToken, certificatePath, diskPath); err != nil {
		return "", err
	}

	log.Logger().Info("Building a new container image...")

	labels, err := vmexport.GetLabelsFromExportSource(virtClient, kind, namespace, name, volumeName)
	if err != nil {
		return "", err
	}

	config := image.DefaultConfig(labels)
	containerImage, err := image.Build(diskPath, config)
	if err != nil {
		return "", err
	}
	digest, err := containerImage.Digest()
	if err != nil {
		return "", err
	}

	log.Logger().Info("Pushing new container image to the container registry...")

	if err := image.Push(containerImage, imageDestination, imagePushTimeout); err != nil {
		return "", err
	}

	log.Logger().Info("Successfully uploaded to the container registry.", zap.String("digest", digest.String()))
	return digest.String(), nil
}

func main() {
	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	if err := cliOptions.Init(); err != nil {
		log.Logger().Error(err.Error())
		os.Exit(constants.InvalidCLIInputExitCode)
	}
	log.Logger().Debug("parsed arguments", zap.Reflect("cliOptions", cliOptions))

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(constants.GenericExitCode)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(constants.GenericExitCode)
	}

	virtClient, err := kubecli.GetKubevirtClient()
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(constants.GenericExitCode)
	}

	imageDigest, err := run(*cliOptions, k8sClient, virtClient)
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(constants.DiskUploaderErrorExitCode)
	}

	results := map[string]string{constants.DigestResultName: imageDigest}

	log.Logger().Debug("recording results", zap.Reflect("results", results))

	if err := res.RecordResults(results); err != nil {
		log.Logger().Error(err.Error())
		os.Exit(constants.WriteResultsExitCode)
	}
}
