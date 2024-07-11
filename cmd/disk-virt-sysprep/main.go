package main

import (
	goarg "github.com/alexflint/go-arg"
	"go.uber.org/zap"

	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/execute"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
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
	executor := execute.NewExecutor(cliOptions, DiskImagePath, "virt-sysprep")

	if err := executor.PrepareGuestFSAppliance(); err != nil {
		exit.ExitOrDieFromError(PrepareGuestFSApplianceFailed, err)
	}

	if err := executor.Execute(); err != nil {
		exit.ExitOrDieFromError(ExecuteFailed, err)
	}
}
