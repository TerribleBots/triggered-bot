package log

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"triggered-bot.log"}
	logger, _ := cfg.Build()
	defer logger.Sync()
	Log = logger.Sugar()
}
