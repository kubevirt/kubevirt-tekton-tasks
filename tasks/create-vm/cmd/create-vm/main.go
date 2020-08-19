package main

import (
	goarg "github.com/alexflint/go-arg"
	. "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/constants"
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils"
	log "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/logger"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/output"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/parse"
	res "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/results"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/vmcreator"
	"net/http"
)

func main() {
	defer utils.HandleExit()

	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	if err := cliOptions.Init(); err != nil {
		utils.ExitOrDieFromError(InvalidNamespacesExitCode, err)
	}

	vmCreator, err := vmcreator.NewVMCreator(cliOptions)

	if err != nil {
		utils.ExitOrDieFromError(GenericExitCode, err)
	}

	if err := vmCreator.CheckVolumesExist(); err != nil {
		utils.ExitFromError(VolumesNotPresentExitCode, err)
	}

	vm, err := vmCreator.CreateVM()

	if err != nil {
		utils.ExitOrDieFromError(CreateVMErrorExitCode, err,
			errors2.IsStatusErrorSoft(err, http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity),
		)
	}

	if err := vmCreator.OwnVolumes(vm); err != nil {
		utils.ExitFromError(OwnVolumesErrorExitCode, err)
	}

	results := map[string]string{
		NameResultName:      vm.Name,
		NamespaceResultName: vm.Namespace,
	}

	if err := res.RecordResults(results); err != nil {
		utils.ExitOrDieFromError(WriteResultsExitCode, err)
	}

	output.PrettyPrint(vm, cliOptions.Output)
}
