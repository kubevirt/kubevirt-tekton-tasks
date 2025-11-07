package vmcreator

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	virtualMachine "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/vm"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/ownerreference"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"k8s.io/client-go/kubernetes"
	"kubevirt.io/client-go/kubecli"

	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	kubevirtv1 "kubevirt.io/api/core/v1"
	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
	virtctlclientconfig "kubevirt.io/kubevirt/pkg/virtctl/clientconfig"
	virtctl "kubevirt.io/kubevirt/pkg/virtctl/create"
	"sigs.k8s.io/yaml"
)

type VMCreator struct {
	targetNamespace        string
	cliOptions             *parse.CLIOptions
	config                 *rest.Config
	virtualMachineProvider virtualMachine.VirtualMachineProvider
	k8sClient              kubernetes.Interface
}

func NewVMCreator(cliOptions *parse.CLIOptions) (*VMCreator, error) {
	log.Logger().Debug("initialized clients and providers")
	targetNS := cliOptions.GetVirtualMachineNamespace()

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// clients
	kubevirtClient, err := kubevirtcliv1.GetKubevirtClientFromRESTConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubevirt client: %v", err.Error())
	}

	k8sclient := kubernetes.NewForConfigOrDie(config)

	virtualMachineProvider := virtualMachine.NewVirtualMachineProvider(kubevirtClient)

	return &VMCreator{
		targetNamespace:        targetNS,
		cliOptions:             cliOptions,
		config:                 config,
		virtualMachineProvider: virtualMachineProvider,
		k8sClient:              k8sclient,
	}, nil
}

func (v *VMCreator) StartVM(namespace, name string) error {
	return v.virtualMachineProvider.Start(namespace, name)
}

func (v *VMCreator) CreateVM() (*kubevirtv1.VirtualMachine, error) {
	switch v.cliOptions.GetCreationMode() {
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

	if v.cliOptions.GetSetOwnerReferenceValue() {
		if err := ownerreference.SetPodOwnerReference(v.k8sClient, &vm); err != nil {
			return nil, err
		}
	}

	return v.virtualMachineProvider.Create(namespace, &vm)
}

func runCommand(params string) ([]byte, error) {
	args := strings.Split(params, " ")
	output := &bytes.Buffer{}

	cmd := virtctl.NewCommand()

	ctx := context.Background()
	cmdContext := virtctlclientconfig.NewContext(ctx, kubecli.DefaultClientConfig(cmd.PersistentFlags()))

	cmd.SetContext(cmdContext)
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
	virtualMachine.AddMetadata(&vm)

	runStrategy := kubevirtv1.VirtualMachineRunStrategy(v.cliOptions.GetRunStrategy())
	if runStrategy != "" {
		vm.Spec.Running = nil
		vm.Spec.RunStrategy = &runStrategy
	}

	if v.cliOptions.GetSetOwnerReferenceValue() {
		if err := ownerreference.SetPodOwnerReference(v.k8sClient, &vm); err != nil {
			return nil, err
		}
	}

	log.Logger().Debug("creating VM", zap.Reflect("vm", vm))
	return v.virtualMachineProvider.Create(v.targetNamespace, &vm)
}
