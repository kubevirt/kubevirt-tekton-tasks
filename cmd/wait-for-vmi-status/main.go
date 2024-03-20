package main

import (
	goarg "github.com/alexflint/go-arg"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/watch"
	"go.uber.org/zap"
	"os"
)

func main() {
	defer exit.HandleExit()

	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	log.Logger().Debug("parsed arguments", zap.Reflect("cliOptions", cliOptions))
	if err := cliOptions.Init(); err != nil {
		exit.ExitOrDieFromError(InvalidArguments, err)
	}

	watchFacade, err := watch.NewWatchFacade(cliOptions)
	if err != nil {
		exit.ExitOrDieFromError(WatchFacadeInitFailed, err)
	}

	success := watchFacade.WaitForVMIConditions()

	if !success {
		os.Exit(FailureConditionFulfilled)
	}
}
