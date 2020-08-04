package vmcreator

import (
	templatev1 "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	"github.com/pkg/errors"
	. "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/datavolume"
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/pvc"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/templates"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/parse"
	virtualMachine "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/vm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
	datavolumeclientv1alpha1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1alpha1"
	"path/filepath"
)

type VMCreator struct {
	targetNamespace        string
	cliParams              *parse.CLIParams
	config                 *rest.Config
	templateProvider       templates.TemplateProvider
	virtualMachineProvider virtualMachine.VirtualMachineProvider
	dataVolumeProvider     datavolume.DataVolumeProvider
	pvcProvider            pvc.PersistentVolumeClaimProvider
}

func getConfig() (*rest.Config, error) {
	if IsEnvVarTrue(OutOfClusterENV) {
		return clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	}
	return rest.InClusterConfig()
}

func NewVMCreator(cliParams *parse.CLIParams) (*VMCreator, error) {
	targetNS := cliParams.VirtualMachineNamespace

	if targetNS == "" {
		activeNamespace, err := GetActiveNamespace()
		if err != nil {
			return nil, errors2.NewMissingRequiredError("%v: %v option is empty", err.Error(), parse.VMNamespaceOptionName)
		}
		targetNS = activeNamespace
	}

	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	// clients
	kubeClient := kubernetes.NewForConfigOrDie(config)
	templateClient := templatev1.NewForConfigOrDie(config)
	cdiClient := datavolumeclientv1alpha1.NewForConfigOrDie(config)
	kubevirtClient, err := kubevirtcliv1.GetKubevirtClientFromRESTConfig(config)
	if err != nil {
		return nil, errors.WithMessage(err, "Cannot create kubevirt client")
	}

	templateProvider := templates.NewTemplateProvider(templateClient)
	virtualMachineProvider := virtualMachine.NewVirtualMachineProvider(kubevirtClient)
	dataVolumeProvider := datavolume.NewDataVolumeProvider(cdiClient)
	pvcProvider := pvc.NewPersistentVolumeClaimProvider(kubeClient.CoreV1())

	return &VMCreator{
		targetNamespace:        targetNS,
		cliParams:              cliParams,
		config:                 config,
		templateProvider:       templateProvider,
		virtualMachineProvider: virtualMachineProvider,
		dataVolumeProvider:     dataVolumeProvider,
		pvcProvider:            pvcProvider,
	}, nil
}

func (v *VMCreator) CreateVM() (*kubevirtv1.VirtualMachine, error) {
	templateNamespace := v.cliParams.TemplateNamespace
	if templateNamespace == "" {
		templateNamespace = v.targetNamespace
	}

	template, err := v.templateProvider.Get(templateNamespace, v.cliParams.TemplateName)
	if err != nil {
		return nil, err
	}

	processedTemplate, err := v.templateProvider.Process(v.targetNamespace, template, v.cliParams.GetTemplateParams())
	if err != nil {
		return nil, err
	}
	vm, err := templates.DecodeVM(processedTemplate)
	if err != nil {
		return nil, err
	}

	vm.Namespace = v.targetNamespace
	virtualMachine.AddMetadata(vm, processedTemplate)
	virtualMachine.AddVolumes(vm, processedTemplate, v.cliParams)

	return v.virtualMachineProvider.Create(v.targetNamespace, vm)
}

func (v *VMCreator) CheckVolumesExist() error {
	_, dvsErr := v.dataVolumeProvider.GetByName(v.targetNamespace, v.cliParams.GetAllDVNames()...)
	_, pvcsErr := v.pvcProvider.GetByName(v.targetNamespace, v.cliParams.GetAllPVCNames()...)

	return errors2.NewMultiError().
		AddC("dvsErr", dvsErr).
		AddC("pvcsErr", pvcsErr).
		AsOptional()
}

func (v *VMCreator) OwnVolumes(vm *kubevirtv1.VirtualMachine) error {
	dvsErr := v.ownDataVolumes(vm)
	pvcsErr := v.ownPersistentVolumeClaims(vm)

	return errors2.NewMultiError().
		AddC("dvsErr", dvsErr).
		AddC("pvcsErr", pvcsErr).
		AsOptional()
}

func (v *VMCreator) ownDataVolumes(vm *kubevirtv1.VirtualMachine) error {
	var multiError errors2.MultiError

	dvs, dvsErr := v.dataVolumeProvider.GetByName(v.targetNamespace, v.cliParams.OwnDataVolumes...)

	for idx, dvName := range v.cliParams.OwnDataVolumes {
		if err := errors2.GetErrorFromMultiError(dvsErr, dvName); err != nil {
			multiError.Add(dvName, err)
			continue
		}

		if _, err := v.dataVolumeProvider.AddOwnerReferences(dvs[idx], virtualMachine.AsVMOwnerReference(vm)); err != nil {
			multiError.Add(dvName, errors.Wrapf(err, "could not add owner reference to %v DataVolume", dvName))
		}

	}

	return multiError.AsOptional()
}

func (v *VMCreator) ownPersistentVolumeClaims(vm *kubevirtv1.VirtualMachine) error {
	var multiError errors2.MultiError

	pvcs, pvcsErr := v.pvcProvider.GetByName(v.targetNamespace, v.cliParams.OwnPersistentVolumeClaims...)

	for idx, pvcName := range v.cliParams.OwnPersistentVolumeClaims {
		if err := errors2.GetErrorFromMultiError(pvcsErr, pvcName); err != nil {
			multiError.Add(pvcName, err)
			continue
		}

		if _, err := v.pvcProvider.AddOwnerReferences(pvcs[idx], virtualMachine.AsVMOwnerReference(vm)); err != nil {
			multiError.Add(pvcName, errors.Wrapf(err, "could not add owner reference to %v PersistentVolumeClaim", pvcName))
		}

	}

	return multiError.AsOptional()
}
