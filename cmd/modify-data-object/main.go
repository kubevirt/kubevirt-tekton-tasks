package main

import (
	"os"

	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-data-object/pkg/constants"
	dataobjectcreator "github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-data-object/pkg/dataobject"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/modify-data-object/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	res "github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/results"
	"go.uber.org/zap"
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

	dataObjectCreator, err := dataobjectcreator.NewDataObjectCreator(cliOptions)
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(DataObjectCreatorErrorCode)
	}

	if cliOptions.GetDeleteObject() {
		if err := dataObjectCreator.DeleteDataObject(); err != nil {
			log.Logger().Error(err.Error())
			os.Exit(DeleteObjectExitCode)
		}
		log.Logger().Debug("Object was deleted")
		return
	}

	newDataObject, err := dataObjectCreator.CreateDataObject()
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(CreateDataObjectErrorCode)
	}

	results := map[string]string{
		NameResultName:      newDataObject.GetName(),
		NamespaceResultName: newDataObject.GetNamespace(),
	}

	log.Logger().Debug("recording results", zap.Reflect("results", results))
	if err = res.RecordResults(results); err != nil {
		log.Logger().Error(err.Error())
		os.Exit(WriteResultsExitCode)
	}

	output.PrettyPrint(newDataObject, cliOptions.Output)
}
