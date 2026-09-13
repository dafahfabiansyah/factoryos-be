package platform

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(env string) {
	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := cfg.Build()
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}

	Log = logger
	zap.ReplaceGlobals(logger)
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

func With(fields ...zap.Field) *zap.Logger {
	if Log == nil {
		return zap.L().With(fields...)
	}
	return Log.With(fields...)
}
