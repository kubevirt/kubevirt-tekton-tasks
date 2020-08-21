package main

import (
	goarg "github.com/alexflint/go-arg"
	. "github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/constants"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/log"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/output"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/utils/parse"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/create-vm/pkg/vmcreator"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	res "github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/results"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	defer exit.HandleExit()

	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	if err := cliOptions.Init(); err != nil {
		exit.ExitOrDieFromError(InvalidNamespacesExitCode, err)
	}
	log.GetLogger().Debug("parsed arguments", zap.Reflect("cliOptions", cliOptions))

	vmCreator, err := vmcreator.NewVMCreator(cliOptions)

	if err != nil {
		exit.ExitOrDieFromError(GenericExitCode, err)
	}

	if err := vmCreator.CheckVolumesExist(); err != nil {
		exit.ExitFromError(VolumesNotPresentExitCode, err)
	}

	vm, err := vmCreator.CreateVM()

	if err != nil {
		exit.ExitOrDieFromError(CreateVMErrorExitCode, err,
			zerrors.IsStatusErrorSoft(err, http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity),
		)
	}

	if err := vmCreator.OwnVolumes(vm); err != nil {
		exit.ExitFromError(OwnVolumesErrorExitCode, err)
	}

	results := map[string]string{
		NameResultName:      vm.Name,
		NamespaceResultName: vm.Namespace,
	}

	log.GetLogger().Debug("recording results", zap.Reflect("results", results))
	if err := res.RecordResults(results); err != nil {
		exit.ExitOrDieFromError(WriteResultsExitCode, err)
	}

	output.PrettyPrint(vm, cliOptions.Output)
}
