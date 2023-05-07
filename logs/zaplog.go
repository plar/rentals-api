package logs

import (
	"go.uber.org/zap"
)

func Init() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	log, _ := cfg.Build()
	return log
}
