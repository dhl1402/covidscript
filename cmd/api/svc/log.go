package svc

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func MakeLogger() *zap.Logger {
	zconf := zap.NewProductionConfig()
	zconf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zconf.DisableStacktrace = true
	logger, err := zconf.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
