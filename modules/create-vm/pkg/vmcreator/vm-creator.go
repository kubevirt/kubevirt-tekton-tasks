package vmcreator

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/templates"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	virtualMachine "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/vm"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zutils"
	templatev1 "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	kubevirtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"
	kubevirtcliv1 "kubevirt.io/client-go/kubecli"
	virtctl "kubevirt.io/kubevirt/pkg/virtctl/create"
	"sigs.k8s.io/yaml"
)

type VMCreator struct {
	targetNamespace        string
	cliOptions             *parse.CLIOptions
	config                 *rest.Config
	templateProvider       templates.TemplateProvider
	virtualMachineProvider virtualMachine.VirtualMachineProvider
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

	var templateProvider templates.TemplateProvider
	virtualMachineProvider := virtualMachine.NewVirtualMachineProvider(kubevirtClient)

	if cliOptions.GetCreationMode() == constants.TemplateCreationMode {
		templateProvider = templates.NewTemplateProvider(templatev1.NewForConfigOrDie(config))
	}

	return &VMCreator{
		targetNamespace:        targetNS,
		cliOptions:             cliOptions,
		config:                 config,
		templateProvider:       templateProvider,
		virtualMachineProvider: virtualMachineProvider,
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
	rootCmd := &cobra.Command{
		Use:           "kubevirt-tekton-tasks-create-vm",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf(cmd.UsageString())
		},
	}
	clientConfig := kubecli.DefaultClientConfig(rootCmd.PersistentFlags())
	cmd := virtctl.NewCommand(clientConfig)
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

	vm.Namespace = v.targetNamespace

	runStrategy := kubevirtv1.VirtualMachineRunStrategy(v.cliOptions.GetRunStrategy())
	if runStrategy != "" {
		vm.Spec.Running = nil
		vm.Spec.RunStrategy = &runStrategy
	}

	log.Logger().Debug("creating VM", zap.Reflect("vm", vm))
	return v.virtualMachineProvider.Create(v.targetNamespace, vm)
}
