package main

import (
	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execattributes"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execute"
	log "github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/parse"
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

	execAttributes := execattributes.NewExecAttributes()

	if err := execAttributes.Init(ConnectionSecretPath); err != nil {
		exit.ExitOrDieFromError(InvalidSecret, err)
	}
	log.GetLogger().Debug("retrieved connection secret exec attributes", zap.Object("execAttributes", execAttributes))

	executor, _ := execute.NewExecutor(cliOptions, execAttributes)

	if err := executor.EnsureVMRunning(); err != nil {
		exit.ExitOrDieFromError(EnsureVMRunningFailed, err)
	}

	if err := executor.SetupConnection(); err != nil {
		exit.ExitOrDieFromError(ExecutorSetupConnectionFailed, err)
	}

	if err := executor.RemoteExecute(); err != nil {
		log.GetLogger().Debug("finished", zap.Reflect("err", err))
		exit.ExitOrDieFromError(RemoteCommandExecutionFailedFromUnknownError, err)
	}
}
