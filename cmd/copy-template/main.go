package main

import (
	"net/http"

	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/constants"
	templatecreator "github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/templates"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/copy-template/pkg/utils/parse"
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

	templateCreator, err := templatecreator.NewTemplateCreator(cliOptions)
	if err != nil {
		exit.ExitOrDieFromError(TemplateCreatorErrorCode, err)
	}

	newTemplate, err := templateCreator.CopyTemplate()
	if err != nil {
		exit.ExitOrDieFromError(CopyTemplateErrorCode, err,
			zerrors.IsStatusError(err, http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity),
		)
	}

	results := map[string]string{
		NameResultName:      newTemplate.Name,
		NamespaceResultName: newTemplate.Namespace,
	}

	log.Logger().Debug("recording results", zap.Reflect("results", results))
	err = res.RecordResults(results)
	if err != nil {
		exit.ExitOrDieFromError(WriteResultsExitCode, err)
	}

	output.PrettyPrint(newTemplate, cliOptions.Output)
}
