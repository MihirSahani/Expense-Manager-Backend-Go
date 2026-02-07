package elogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func NewLogger() Logger {
	cfg := zap.NewDevelopmentConfig()                                // console, human-readable [web:21]
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // colored levels [web:16]
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return &zapLogger{logger: logger}
}