package log

import (
	"go.uber.org/zap"
)

func Root() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func Infof(template string, args ...interface{}) {
	logger := Root()
	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Infof(template, args)
}
