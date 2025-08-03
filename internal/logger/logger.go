package logger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type LoggerService struct {
	*zerolog.Logger
}

func New(format, level string) *LoggerService {
	return NewWithWriter(format, level, os.Stderr)
}

func NewWithWriter(format, level string, writer io.Writer) *LoggerService {
	format = strings.ToLower(format)
	level = strings.ToLower(level)

	var logWriter io.Writer
	if format == "json" {
		logWriter = writer
	} else {
		logWriter = zerolog.ConsoleWriter{Out: writer}
	}

	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		parsedLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(logWriter).
		Level(parsedLevel).
		With().
		Timestamp().
		Logger()

	return &LoggerService{
		Logger: &logger,
	}
}

func (ls *LoggerService) WithField(key string, value interface{}) *LoggerService {
	newLogger := ls.Logger.With().Interface(key, value).Logger()
	return &LoggerService{Logger: &newLogger}
}

func (ls *LoggerService) WithFields(fields map[string]interface{}) *LoggerService {
	ctx := ls.Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	newLogger := ctx.Logger()
	return &LoggerService{Logger: &newLogger}
}
