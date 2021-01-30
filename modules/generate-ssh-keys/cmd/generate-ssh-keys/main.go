package main

import (
	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/generate"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/secretcreator"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	res "github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/results"
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

	keys, err := generate.GenerateSshKeys(*cliOptions)
	if err != nil {
		exit.ExitOrDieFromError(SshKeysGenerationFailed, err)
	}

	secretCreator, err := secretcreator.NewSecretCreator(cliOptions, *keys)
	if err != nil {
		exit.ExitOrDieFromError(SecretCreatorInitFailed, err)
	}

	err = secretCreator.CheckPrivateKeySecretExistence()
	if err != nil {
		exit.ExitOrDieFromError(PrivateKeyAlreadyExists, err)
	}

	publicKeySecret, err := secretCreator.AppendPublicKeySecret()
	if err != nil {
		exit.ExitOrDieFromError(PublicKeySecretCreationFailed, err)
	}

	privateKeySecret, err := secretCreator.CreatePrivateKeySecret()
	if err != nil {
		exit.ExitOrDieFromError(PrivateKeySecretCreationFailed, err)
	}

	results := map[string]string{
		Results.PublicKeySecretName:       publicKeySecret.Name,
		Results.PublicKeySecretNamespace:  publicKeySecret.Namespace,
		Results.PrivateKeySecretName:      privateKeySecret.Name,
		Results.PrivateKeySecretNamespace: privateKeySecret.Namespace,
	}

	log.GetLogger().Debug("recording results", zap.Reflect("results", results))
	if err := res.RecordResults(results); err != nil {
		exit.ExitOrDieFromError(WriteResultsExitCode, err)
	}
}
