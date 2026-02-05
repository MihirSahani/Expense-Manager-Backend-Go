package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	*zap.Logger
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}


func NewLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()               // console, human-readable [web:21]
    cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // colored levels [web:16]
    logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}