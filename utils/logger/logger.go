package logger

import (
	"sync"
)

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

type LoggerType int

const (
	Zap LoggerType = iota
	Std
	Slog
)

var (
	once   sync.Once
	logger Logger
)

func InitLogger(logType LoggerType) Logger {
	once.Do(func() {
		switch logType {
		case Zap:
			logger = NewZap()
		// case Slog:
		// 	logger = nil
		default:
			logger = NewStd()
		}
	})

	return logger
}
