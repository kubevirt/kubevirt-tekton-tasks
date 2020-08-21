package log

import (
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/suomiy/kubevirt-tekton-tasks/modules/shared/pkg/zconstants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger(level zapcore.Level) *zap.Logger {
	if logger != nil {
		return logger
	}

	var err error
	var config zap.Config

	// set opinionated presets
	if env.IsEnvVarTrue(zconstants.OutOfClusterENV) {
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
