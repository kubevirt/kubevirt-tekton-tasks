package main

import (
	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/execute"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt-customize/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"go.uber.org/zap"
)

func main() {
	defer exit.HandleExit()

	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	log.GetLogger().Debug("parsed arguments", zap.Reflect("cliOptions", cliOptions))
	if err := cliOptions.Init(); err != nil {
		exit.ExitOrDieFromError(InvalidArguments, err)
	}
	executor := execute.NewExecutor(cliOptions, DiskImagePath)

	if err := executor.PrepareGuestFSAppliance(); err != nil {
		exit.ExitOrDieFromError(PrepareGuestFSApplianceFailed, err)
	}

	if err := executor.Execute(); err != nil {
		exit.ExitOrDieFromError(ExecuteFailed, err)
	}
}
