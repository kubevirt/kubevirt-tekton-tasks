package main

import (
	goarg "github.com/alexflint/go-arg"
	. "github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/constants"
	log "github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/logger"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/execute-in-vm/pkg/utils/parse"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/exit"
)

func main() {
	defer exit.HandleExit()

	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	if err := cliOptions.Init(); err != nil {
		exit.ExitOrDieFromError(InvalidArguments, err)
	}

	// ssh -oStrictHostKeyChecking=accept-new -i /tmp/id_rsa fedora@10.116.0.104
	// -oStrictHostKeyChecking=no
}
