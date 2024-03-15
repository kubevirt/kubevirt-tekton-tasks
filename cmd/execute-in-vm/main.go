package main

import (
	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execute"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils"
	log "github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
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

	executor, executorErr := execute.NewExecutor(cliOptions, ConnectionSecretPath)
	if executorErr != nil {
		exit.ExitOrDieFromError(ExecutorInitialization, executorErr)
	}

	multiError := zerrors.NewMultiError()
	var exitError *exit.Exit

	registerError := func(name string, err error) {
		if err != nil {
			if exitErr, ok := err.(exit.Exit); ok {
				exitError = &exitErr
			} else if err == wait.ErrWaitTimeout {
				exitError = &exit.Exit{
					Code: CommandTimeout,
					Msg:  "command timed out",
					Soft: true,
				}
			} else {
				multiError.Add(name, err)
			}
		}

	}

	if cliOptions.GetScript() != "" {
		runWithTimeout := utils.WithTimeout(cliOptions.GetScriptTimeout())

		runWithTimeout(func(timeout time.Duration, finished bool) {
			if multiError.IsEmpty() && !finished {
				err := executor.EnsureVMRunning(timeout)
				registerError("EnsureVMRunning", err)
			}
		})

		runWithTimeout(func(timeout time.Duration, finished bool) {
			if multiError.IsEmpty() && !finished {
				err := executor.SetupConnection(timeout)
				registerError("SetupConnection", err)
			}
		})

		runWithTimeout(func(timeout time.Duration, finished bool) {
			if multiError.IsEmpty() {
				if !finished {
					err := executor.RemoteExecute(timeout)
					registerError("RemoteExecute", err)
				} else {
					registerError("RemoteExecute", wait.ErrWaitTimeout)
				}

			}
		})

	}

	if cliOptions.ShouldStop() {
		if err := executor.EnsureVMStopped(); err != nil {
			multiError.Add("VM Stop", err)
		}
	}

	if cliOptions.ShouldDelete() {
		if err := executor.EnsureVMDeleted(); err != nil {
			multiError.Add("VM Delete", err)
		}
	}

	if !multiError.IsEmpty() {
		if exitError != nil {
			multiError.Add("command exit", *exitError)
		}
		log.Logger().Debug("finished", zap.String("errMsg", multiError.Error()))
		exit.ExitOrDieFromError(ExecutorActionsFailed, multiError)
	}

	if exitError != nil {
		log.Logger().Debug("finished", zap.Reflect("err", exitError))
		exit.ExitOrDieFromError(exitError.Code, exitError)
	}
}
