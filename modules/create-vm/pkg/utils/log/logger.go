package log

import (
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
	config = zap.NewProductionConfig()

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
