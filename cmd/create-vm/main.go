package main

import (
	"os"

	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-vm/pkg/vmcreator"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	res "github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/results"
	"go.uber.org/zap"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

func main() {
	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	if err := cliOptions.Init(); err != nil {
		log.Logger().Error(err.Error())
		os.Exit(InvalidCLIInputExitCode)
	}
	log.Logger().Debug("parsed arguments", zap.Reflect("cliOptions", cliOptions))

	vmCreator, err := vmcreator.NewVMCreator(cliOptions)
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(GenericExitCode)
	}

	vm, err := vmCreator.CreateVM()
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(CreateVMErrorExitCode)
	}

	if cliOptions.GetStartVMFlag() &&
		(vm.Spec.RunStrategy == nil || *vm.Spec.RunStrategy != kubevirtv1.RunStrategyAlways) &&
		(vm.Spec.Running == nil || !*vm.Spec.Running) {
		if err := vmCreator.StartVM(vm.Namespace, vm.Name); err != nil {
			log.Logger().Error(err.Error())
			os.Exit(StartVMErrorExitCode)
		}
	}

	results := map[string]string{
		NameResultName:      vm.Name,
		NamespaceResultName: vm.Namespace,
	}

	log.Logger().Debug("recording results", zap.Reflect("results", results))
	if err := res.RecordResults(results); err != nil {
		log.Logger().Error(err.Error())
		os.Exit(WriteResultsExitCode)
	}

	output.PrettyPrint(vm, cliOptions.Output)
}
