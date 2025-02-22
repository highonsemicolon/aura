package logger

import "log"

type StdLogger struct{}

func NewStd() *StdLogger {
	return &StdLogger{}
}

func (l *StdLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Println("[INFO]", msg, keysAndValues)
}

func (l *StdLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Println("[ERROR]", msg, keysAndValues)
}
