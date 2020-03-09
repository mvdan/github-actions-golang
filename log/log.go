package log

import (
	"go.uber.org/zap"
)

func Infof(template string, args ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Infof(template, args)
}
