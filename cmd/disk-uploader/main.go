package main

import (
	"os"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/certificate"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/disk"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/image"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/secrets"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-uploader/pkg/vmexport"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	goarg "github.com/alexflint/go-arg"
	"go.uber.org/zap"
	kubecli "kubevirt.io/client-go/kubecli"
)

const (
	genericExitCode           = 1
	invalidCLIInputExitCode   = 2
	diskUploaderErrorExitCode = 3

	diskPath            string = "./tmp/disk.qcow2"
	certificatePath     string = "./tmp/tls.crt"
	kvExportTokenHeader string = "x-kubevirt-export-token"
)

func run(opts parse.CLIOptions, k8sClient kubernetes.Interface, virtClient kubecli.KubevirtClient) error {
	kind := opts.GetExportSourceKind()
	name := opts.GetExportSourceName()
	namespace := opts.GetExportSourceNamespace()
	volumeName := opts.GetVolumeName()
	imageDestination := opts.GetImageDestination()
	imagePushTimeout := opts.GetPushTimeout()

	log.Logger().Info("Creating a new Secret object...", zap.String("namespace", namespace), zap.String("name", name))

	if err := secrets.CreateVirtualMachineExportSecret(k8sClient, namespace, name); err != nil {
		return err
	}

	log.Logger().Info("Creating a new VirtualMachineExport object...", zap.String("namespace", namespace), zap.String("name", name))

	if err := vmexport.CreateVirtualMachineExport(virtClient, kind, namespace, name); err != nil {
		return err
	}

	log.Logger().Info("Waiting for VirtualMachineExport status to be ready...")

	if err := vmexport.WaitUntilVirtualMachineExportReady(virtClient, namespace, name); err != nil {
		return err
	}

	log.Logger().Info("Getting raw disk URL from the VirtualMachineExport object status...")

	rawDiskUrl, err := vmexport.GetRawDiskUrlFromVolumes(virtClient, namespace, name, volumeName)
	if err != nil {
		return err
	}

	log.Logger().Info("Creating TLS certificate file from the VirtualMachineExport object status...")

	certificateData, err := certificate.GetCertificateFromVirtualMachineExport(virtClient, namespace, name)
	if err != nil {
		return err
	}

	if err := os.WriteFile(certificatePath, []byte(certificateData), 0644); err != nil {
		return err
	}

	log.Logger().Info("Getting export token from the Secret object...")

	kvExportToken, err := secrets.GetTokenFromVirtualMachineExportSecret(virtClient, namespace, name)
	if err != nil {
		return err
	}

	log.Logger().Info("Downloading disk image from the VirtualMachineExport server...")

	if err := disk.DownloadDiskImageFromURL(rawDiskUrl, kvExportTokenHeader, kvExportToken, certificatePath, diskPath); err != nil {
		return err
	}

	log.Logger().Info("Building a new container image...")

	containerImage, err := image.Build(diskPath)
	if err != nil {
		return err
	}

	log.Logger().Info("Pushing new container image to the container registry...")

	if err := image.Push(containerImage, imageDestination, imagePushTimeout); err != nil {
		return err
	}

	log.Logger().Info("Successfully uploaded to the container registry.")
	return nil
}

func main() {
	defer exit.HandleExit()

	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	if err := cliOptions.Init(); err != nil {
		exit.ExitOrDieFromError(invalidCLIInputExitCode, err)
	}
	log.Logger().Debug("parsed arguments", zap.Reflect("cliOptions", cliOptions))

	config, err := rest.InClusterConfig()
	if err != nil {
		exit.ExitOrDieFromError(genericExitCode, err)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		exit.ExitOrDieFromError(genericExitCode, err)
	}

	virtClient, err := kubecli.GetKubevirtClient()
	if err != nil {
		exit.ExitOrDieFromError(genericExitCode, err)
	}

	if err := run(*cliOptions, k8sClient, virtClient); err != nil {
		exit.ExitOrDieFromError(diskUploaderErrorExitCode, err)
	}
}
