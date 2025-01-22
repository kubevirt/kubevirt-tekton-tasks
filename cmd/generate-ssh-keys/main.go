package main

import (
	"os"

	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/generate"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/secret"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/parse"
	res "github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/results"
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

	keys, err := generate.GenerateSshKeys(*cliOptions)
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(SshKeysGenerationFailed)
	}

	secretFacade, err := secret.NewSecretFacade(cliOptions, *keys)
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(SecretFacadeInitFailed)
	}

	err = secretFacade.CheckPrivateKeySecretExistence()
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(PrivateKeyAlreadyExists)
	}

	publicKeySecret, err := secretFacade.GetPublicKeySecret()
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(PublicKeySecretFetchFailed)
	}
	isAppendingPublicKey := publicKeySecret != nil

	if isAppendingPublicKey {
		publicKeySecret, err = secretFacade.AppendPublicKeySecret(publicKeySecret)
	} else {
		publicKeySecret, err = secretFacade.CreatePublicKeySecret()
	}

	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(PublicKeySecretCreationFailed)
	}

	cleanupPublicKey := func() {
		if !isAppendingPublicKey {
			_ = secretFacade.DeleteSecret(publicKeySecret)
		}
	}

	privateKeySecret, err := secretFacade.CreatePrivateKeySecret()
	if err != nil {
		defer cleanupPublicKey()
		log.Logger().Error(err.Error())
		os.Exit(PrivateKeySecretCreationFailed)
	}

	cleanupPrivateKey := func() {
		_ = secretFacade.DeleteSecret(privateKeySecret)
	}

	results := map[string]string{
		Results.PublicKeySecretName:       publicKeySecret.Name,
		Results.PublicKeySecretNamespace:  publicKeySecret.Namespace,
		Results.PrivateKeySecretName:      privateKeySecret.Name,
		Results.PrivateKeySecretNamespace: privateKeySecret.Namespace,
	}

	log.Logger().Debug("recording results", zap.Reflect("results", results))
	if err := res.RecordResults(results); err != nil {
		defer func() {
			defer cleanupPublicKey()
			defer cleanupPrivateKey()
		}()
		log.Logger().Error(err.Error())
		os.Exit(WriteResultsExitCode)
	}
}
