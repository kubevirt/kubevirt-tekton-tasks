package main

import (
	goarg "github.com/alexflint/go-arg"
	. "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/constants"
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/output"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/parse"
	res "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/results"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/vmcreator"
	"net/http"
)

func main() {
	cliParams := &parse.CLIParams{}
	goarg.MustParse(cliParams)

	vmCreator, err := vmcreator.NewVMCreator(cliParams)

	if err != nil {
		utils.ErrorExitOrDie(GenericExitCode, err)
	}

	if err := vmCreator.CheckVolumesExist(); err != nil {
		utils.ErrorExit(VolumesNotPresentExitCode, err)
	}

	vm, err := vmCreator.CreateVM()

	if err != nil {
		utils.ErrorExitOrDie(CreateVMErrorExitCode, err,
			errors2.IsStatusErrorSoft(err, http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity),
		)
	}

	if err := vmCreator.OwnVolumes(vm); err != nil {
		utils.ErrorExit(OwnVolumesErrorExitCode, err)
	}

	results := map[string]string{
		NameResultName:      vm.Name,
		NamespaceResultName: vm.Namespace,
	}

	if err := res.WriteResults(results); err != nil {
		utils.ErrorExitOrDie(WriteResultsExitCode, err)
	}

	output.PrettyPrint(vm, cliParams.Output)
}
