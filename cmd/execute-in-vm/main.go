package main

import (
	"fmt"
	"os"
	"time"

	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execute"
	log "github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/parse"
	"go.uber.org/zap"
)

func executeScript(executor *execute.Executor, scriptTimeout time.Duration) error {
	deadline := time.Now().Add(scriptTimeout)

	if err := executor.EnsureVMRunning(scriptTimeout); err != nil {
		return err
	}

	remaining := time.Until(deadline)
	if scriptTimeout > 0 && remaining <= 0 {
		return fmt.Errorf("timed out waiting for VM to be running")
	}

	if err := executor.SetupConnection(remaining); err != nil {
		return err
	}

	remaining = time.Until(deadline)
	if scriptTimeout > 0 && remaining <= 0 {
		return fmt.Errorf("timed out waiting for SSH connection")
	}

	return executor.RemoteExecute(remaining)
}

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

	executor, executorErr := execute.NewExecutor(cliOptions, ConnectionSecretPath)
	if executorErr != nil {
		log.Logger().Error(executorErr.Error())
		os.Exit(ExecutorInitialization)
	}

	if cliOptions.GetScript() != "" {
		if err := executeScript(executor, cliOptions.GetScriptTimeout()); err != nil {
			log.Logger().Error(err.Error())
			os.Exit(ExecutorActionsFailed)
		}
	}

	if cliOptions.ShouldStop() {
		if err := executor.EnsureVMStopped(); err != nil {
			log.Logger().Error(err.Error())
			os.Exit(ExecutorActionsFailed)
		}
	}

	if cliOptions.ShouldDelete() {
		if err := executor.EnsureVMDeleted(); err != nil {
			log.Logger().Error(err.Error())
			os.Exit(ExecutorActionsFailed)
		}
	}
}
