package main

import (
	"net/http"

	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-data-object/pkg/constants"
	dataobjectcreator "github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-data-object/pkg/dataobject"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-data-object/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	res "github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/results"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap"
)

func main() {
	defer exit.HandleExit()

	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	err := cliOptions.Init()
	if err != nil {
		exit.ExitOrDieFromError(InvalidCLIInputExitCode, err)
	}

	log.Logger().Debug("parsed arguments", zap.Reflect("cliOptions", cliOptions))

	dataObjectCreator, err := dataobjectcreator.NewDataObjectCreator(cliOptions)
	if err != nil {
		exit.ExitOrDieFromError(DataObjectCreatorErrorCode, err)
	}

	if cliOptions.GetDeleteObject() {
		err := dataObjectCreator.DeleteDataObject()
		if err != nil {
			exit.ExitOrDieFromError(DeleteObjectExitCode, err)
		}
		log.Logger().Debug("Object was deleted")
		return
	}

	newDataObject, err := dataObjectCreator.CreateDataObject()
	if err != nil {
		exit.ExitOrDieFromError(CreateDataObjectErrorCode, err,
			zerrors.IsStatusError(err, http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity),
		)
	}

	results := map[string]string{
		NameResultName:      newDataObject.GetName(),
		NamespaceResultName: newDataObject.GetNamespace(),
	}

	log.Logger().Debug("recording results", zap.Reflect("results", results))
	err = res.RecordResults(results)
	if err != nil {
		exit.ExitOrDieFromError(WriteResultsExitCode, err)
	}

	output.PrettyPrint(newDataObject, cliOptions.Output)
}
