package vmcreator

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/datavolume"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/pvc"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates/validations"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	virtualMachine "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/vm"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	templatev1 "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	kubevirtv1 "kubevirt.io/api/core/v1"
	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
	datavolumeclientv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"
	virtctl "kubevirt.io/kubevirt/pkg/virtctl/create"
	"sigs.k8s.io/yaml"
)

type VMCreator struct {
	targetNamespace        string
	cliOptions             *parse.CLIOptions
	config                 *rest.Config
	templateProvider       templates.TemplateProvider
	virtualMachineProvider virtualMachine.VirtualMachineProvider
	dataVolumeProvider     datavolume.DataVolumeProvider
	pvcProvider            pvc.PersistentVolumeClaimProvider
}

func NewVMCreator(cliOptions *parse.CLIOptions) (*VMCreator, error) {
	log.Logger().Debug("initialized clients and providers")
	targetNS := cliOptions.GetVirtualMachineNamespace()

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// clients
	kubeClient := kubernetes.NewForConfigOrDie(config)
	cdiClient := datavolumeclientv1beta1.NewForConfigOrDie(config)
	kubevirtClient, err := kubevirtcliv1.GetKubevirtClientFromRESTConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubevirt client: %v", err.Error())
	}

	var templateProvider templates.TemplateProvider
	virtualMachineProvider := virtualMachine.NewVirtualMachineProvider(kubevirtClient)
	dataVolumeProvider := datavolume.NewDataVolumeProvider(cdiClient)
	pvcProvider := pvc.NewPersistentVolumeClaimProvider(kubeClient.CoreV1())

	if cliOptions.GetCreationMode() == constants.TemplateCreationMode {
		templateProvider = templates.NewTemplateProvider(templatev1.NewForConfigOrDie(config))
	}

	return &VMCreator{
		targetNamespace:        targetNS,
		cliOptions:             cliOptions,
		config:                 config,
		templateProvider:       templateProvider,
		virtualMachineProvider: virtualMachineProvider,
		dataVolumeProvider:     dataVolumeProvider,
		pvcProvider:            pvcProvider,
	}, nil
}

func (v *VMCreator) StartVM(namespace, name string) error {
	return v.virtualMachineProvider.Start(namespace, name)
}

func (v *VMCreator) CreateVM() (*kubevirtv1.VirtualMachine, error) {
	switch v.cliOptions.GetCreationMode() {
	case constants.TemplateCreationMode:
		return v.createVMFromTemplate()
	case constants.VMManifestCreationMode:
		return v.createVMFromManifest()
	case constants.VirtctlCreatingMode:
		return v.createVMVirtctl()
	}
	return nil, zerrors.NewMissingRequiredError("unknown creation mode: %v", v.cliOptions.GetCreationMode())
}

func (v *VMCreator) createVMVirtctl() (*kubevirtv1.VirtualMachine, error) {
	var vm kubevirtv1.VirtualMachine

	output, err := runCommand(v.cliOptions.Virtctl)
	if err != nil {
		return nil, zerrors.NewSoftError("failed to execute command: %v", err.Error())
	}

	if err := yaml.Unmarshal(output, &vm); err != nil {
		return nil, zerrors.NewSoftError("could not read from virtctl output: %v", err.Error())
	}

	namespace := v.targetNamespace
	if namespace == "" {
		if namespace, err = env.GetActiveNamespace(); err != nil {
			return nil, zerrors.NewMissingRequiredError("can't get active namespace: %v", err.Error())
		}
	}

	return v.virtualMachineProvider.Create(namespace, &vm)
}

func runCommand(params string) ([]byte, error) {
	args := strings.Split(params, " ")
	output := &bytes.Buffer{}
	cmd := virtctl.NewCommand()
	cmd.SetArgs(append([]string{"vm"}, args...))
	cmd.SetOut(output)
	err := cmd.Execute()

	return output.Bytes(), err
}

func (v *VMCreator) createVMFromManifest() (*kubevirtv1.VirtualMachine, error) {
	var vm kubevirtv1.VirtualMachine

	if err := yaml.Unmarshal([]byte(v.cliOptions.VirtualMachineManifest), &vm); err != nil {
		return nil, zerrors.NewSoftError("could not read VM manifest: %v", err.Error())
	}

	vm.Namespace = v.targetNamespace
	virtualMachine.AddMetadata(&vm, nil)

	templateValidations := validations.NewTemplateValidations(nil) // fallback to defaults
	virtualMachine.AddVolumes(&vm, templateValidations, v.cliOptions)

	runStrategy := kubevirtv1.VirtualMachineRunStrategy(v.cliOptions.GetRunStrategy())
	if runStrategy != "" {
		vm.Spec.Running = nil
		vm.Spec.RunStrategy = &runStrategy
	}

	log.Logger().Debug("creating VM", zap.Reflect("vm", vm))
	return v.virtualMachineProvider.Create(v.targetNamespace, &vm)
}

func (v *VMCreator) createVMFromTemplate() (*kubevirtv1.VirtualMachine, error) {
	log.Logger().Debug("retrieving template", zap.String("name", v.cliOptions.TemplateName), zap.String("namespace", v.cliOptions.GetTemplateNamespace()))
	template, err := v.templateProvider.Get(v.cliOptions.GetTemplateNamespace(), v.cliOptions.TemplateName)
	if err != nil {
		return nil, err
	}

	log.Logger().Debug("processing template", zap.String("name", v.cliOptions.TemplateName), zap.String("namespace", v.cliOptions.GetTemplateNamespace()))
	processedTemplate, err := v.templateProvider.Process(v.targetNamespace, template, v.cliOptions.GetTemplateParams())
	if err != nil {
		return nil, err
	}
	vm, _, err := zutils.DecodeVM(processedTemplate)
	if err != nil {
		return nil, err
	}

	templateValidations, err := templates.GetTemplateValidations(processedTemplate)
	if err != nil {
		log.Logger().Warn("could not parse template validations", zap.Error(err))
		templateValidations = validations.NewTemplateValidations(nil) // fallback to defaults
	}
	if templateValidations.IsEmpty() {
		log.Logger().Debug("template validations are empty: falling back to defaults")
	}

	vm.Namespace = v.targetNamespace

	virtualMachine.AddMetadata(vm, processedTemplate)
	virtualMachine.AddVolumes(vm, templateValidations, v.cliOptions)

	runStrategy := kubevirtv1.VirtualMachineRunStrategy(v.cliOptions.GetRunStrategy())
	if runStrategy != "" {
		vm.Spec.Running = nil
		vm.Spec.RunStrategy = &runStrategy
	}

	log.Logger().Debug("creating VM", zap.Reflect("vm", vm))
	return v.virtualMachineProvider.Create(v.targetNamespace, vm)
}

func (v *VMCreator) CheckVolumesExist() error {
	allDVs := zutils.ConcatStringSlices(v.cliOptions.GetOwnDVNames(), v.cliOptions.GetDVNames())
	allPVCs := zutils.ConcatStringSlices(v.cliOptions.GetOwnPVCNames(), v.cliOptions.GetPVCNames())

	log.Logger().Debug("asserting additional volumes exist", zap.Strings("additional-dvs", allDVs), zap.Strings("additional-pvcs", allPVCs))
	_, notFoundDVs, dvsErr := v.dataVolumeProvider.GetByName(v.targetNamespace, allDVs...)

	for dv, _ := range notFoundDVs {
		allPVCs = append(allPVCs, dv)
	}

	_, pvcsErr := v.pvcProvider.GetByName(v.targetNamespace, allPVCs...)

	return zerrors.NewMultiError().
		AddC("dvsErr", dvsErr).
		AddC("pvcsErr", pvcsErr).
		AsOptional()
}

func (v *VMCreator) OwnVolumes(vm *kubevirtv1.VirtualMachine) error {
	dvsErr := v.ownDataVolumes(vm)
	pvcsErr := v.ownPersistentVolumeClaims(vm)

	return zerrors.NewMultiError().
		AddC("dvsErr", dvsErr).
		AddC("pvcsErr", pvcsErr).
		AsOptional()
}

func (v *VMCreator) ownDataVolumes(vm *kubevirtv1.VirtualMachine) error {
	ownDVs := v.cliOptions.GetOwnDVNames()
	log.Logger().Debug("taking ownership of DataVolumes", zap.Strings("own-dvs", ownDVs))
	var multiError zerrors.MultiError

	dvs, notFoundDVs, dvsErr := v.dataVolumeProvider.GetByName(v.targetNamespace, ownDVs...)

	for idx, dvName := range ownDVs {
		if _, ok := notFoundDVs[dvName]; ok {
			// DV not found, nothing to do
			continue
		}

		if err := zerrors.GetErrorFromMultiError(dvsErr, dvName); err != nil {
			multiError.Add(dvName, err)
			continue
		}

		if _, err := v.dataVolumeProvider.AddOwnerReferences(dvs[idx], virtualMachine.AsVMOwnerReference(vm)); err != nil {
			multiError.Add(dvName, fmt.Errorf("could not add owner reference to %v DataVolume: %v", dvName, err.Error()))
		}

	}

	return multiError.AsOptional()
}

func (v *VMCreator) ownPersistentVolumeClaims(vm *kubevirtv1.VirtualMachine) error {
	ownPVCs := v.cliOptions.GetOwnPVCNames()
	log.Logger().Debug("taking ownership of PersistentVolumeClaims", zap.Strings("own-pvcs", ownPVCs))
	var multiError zerrors.MultiError

	pvcs, pvcsErr := v.pvcProvider.GetByName(v.targetNamespace, ownPVCs...)

	for idx, pvcName := range ownPVCs {
		if err := zerrors.GetErrorFromMultiError(pvcsErr, pvcName); err != nil {
			multiError.Add(pvcName, err)
			continue
		}

		if _, err := v.pvcProvider.AddOwnerReferences(pvcs[idx], virtualMachine.AsVMOwnerReference(vm)); err != nil {
			multiError.Add(pvcName, fmt.Errorf("could not add owner reference to %v PersistentVolumeClaim: %v", pvcName, err.Error()))
		}
	}

	return multiError.AsOptional()
}
