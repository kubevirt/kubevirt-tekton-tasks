package main

import (
	"os"

	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/wait-for-vmi-status/pkg/watch"
	"go.uber.org/zap"
)

func main() {
	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	log.Logger().Debug("parsed arguments", zap.Reflect("cliOptions", cliOptions))
	if err := cliOptions.Init(); err != nil {
		log.Logger().Error(err.Error())
		os.Exit(InvalidArguments)
	}

	watchFacade, err := watch.NewWatchFacade(cliOptions)
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(WatchFacadeInitFailed)
	}

	success := watchFacade.WaitForVMIConditions()

	if !success {
		os.Exit(FailureConditionFulfilled)
	}
}
