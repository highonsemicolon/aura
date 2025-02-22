package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.Logger
}

func NewZap() *ZapLogger {
	logger, _ := zap.NewProduction()
	return &ZapLogger{logger: logger}
}

func (l *ZapLogger) Info(msg string, fields ...interface{}) {
	l.logger.Sugar().Infow(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...interface{}) {
	l.logger.Sugar().Errorw(msg, fields...)
}
