package logger

import (
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger(level zapcore.Level) *zap.Logger {
	var err error
	var config zap.Config

	// set opinionated presets
	if constants.IsEnvVarTrue(constants.OutOfClusterENV) {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.Level.SetLevel(level)

	logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func GetLogger() *zap.Logger {
	return logger
}
