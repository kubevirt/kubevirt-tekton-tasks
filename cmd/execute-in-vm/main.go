package main

import (
	"os"
	"time"

	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/execute"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils"
	log "github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/parse"
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

	executor, executorErr := execute.NewExecutor(cliOptions, ConnectionSecretPath)
	if executorErr != nil {
		log.Logger().Error(executorErr.Error())
		os.Exit(ExecutorInitialization)
	}

	if cliOptions.GetScript() != "" {
		runWithTimeout := utils.WithTimeout(cliOptions.GetScriptTimeout())

		runWithTimeout(func(timeout time.Duration, finished bool) {
			if !finished {
				if err := executor.EnsureVMRunning(timeout); err != nil {
					log.Logger().Error(err.Error())
					os.Exit(ExecutorActionsFailed)
				}
			}
		})

		runWithTimeout(func(timeout time.Duration, finished bool) {
			if !finished {
				if err := executor.SetupConnection(timeout); err != nil {
					log.Logger().Error(err.Error())
					os.Exit(ExecutorActionsFailed)
				}
			}
		})

		runWithTimeout(func(timeout time.Duration, finished bool) {
			if !finished {
				if err := executor.RemoteExecute(timeout); err != nil {
					log.Logger().Error(err.Error())
					os.Exit(ExecutorActionsFailed)
				}
			}
		})

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
