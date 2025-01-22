package main

import (
	"os"

	goarg "github.com/alexflint/go-arg"
	"go.uber.org/zap"

	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/execute"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/disk-virt/pkg/utils/parse"
)

func main() {
	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	log.GetLogger().Debug("parsed arguments", zap.Reflect("cliOptions", cliOptions))
	if err := cliOptions.Init(); err != nil {
		log.GetLogger().Error(err.Error())
		os.Exit(InvalidArguments)
	}
	executor := execute.NewExecutor(cliOptions, DiskImagePath, "virt-customize")

	if err := executor.PrepareGuestFSAppliance(); err != nil {
		log.GetLogger().Error(err.Error())
		os.Exit(PrepareGuestFSApplianceFailed)
	}

	if err := executor.Execute(); err != nil {
		log.GetLogger().Error(err.Error())
		os.Exit(ExecuteFailed)
	}
}
