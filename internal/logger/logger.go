package logger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string, errs ...error)
	Error(msg string, errs ...error)
	Fatal(msg string, errs ...error)

	DebugF(format string, args ...any)
	InfoF(format string, args ...any)
	WarnF(format string, args ...any)
	ErrorF(format string, args ...any)
	FatalF(format string, args ...any)

	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
}

type zerologAdapter struct {
	logger *zerolog.Logger
}

func NewZerologAdapter(format, level string) Logger {
	writer := os.Stdout
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

	return &zerologAdapter{
		logger: &logger,
	}
}

func (z *zerologAdapter) Debug(msg string) {
	z.logger.Debug().Msg(msg)
}

func (z *zerologAdapter) Info(msg string) {
	z.logger.Info().Msg(msg)
}

func (z *zerologAdapter) Warn(msg string, errs ...error) {
	events := z.logger.Warn()
	if len(errs) > 0 {
		events = events.Errs("errors", errs)
	}
	events.Msg(msg)
}

func (z *zerologAdapter) Error(msg string, errs ...error) {
	events := z.logger.Error()
	if len(errs) > 0 {
		events = events.Errs("errors", errs)
	}
	events.Msg(msg)
}

func (z *zerologAdapter) Fatal(msg string, errs ...error) {
	event := z.logger.Fatal()

	if len(errs) > 0 {
		event = event.Errs("errors", errs)
	}

	event.Msg(msg)
}

func (z *zerologAdapter) DebugF(format string, args ...any) {
	z.logger.Debug().Msgf(format, args...)
}
func (z *zerologAdapter) InfoF(format string, args ...any) {
	z.logger.Info().Msgf(format, args...)
}
func (z *zerologAdapter) WarnF(format string, args ...any) {
	z.logger.Warn().Msgf(format, args...)
}
func (z *zerologAdapter) ErrorF(format string, args ...any) {
	z.logger.Error().Msgf(format, args...)
}
func (z *zerologAdapter) FatalF(format string, args ...any) {
	z.logger.Fatal().Msgf(format, args...)
}

func (z *zerologAdapter) WithField(key string, value any) Logger {
	newLogger := z.logger.With().Interface(key, value).Logger()
	return &zerologAdapter{logger: &newLogger}
}

func (z *zerologAdapter) WithFields(fields map[string]any) Logger {
	ctx := z.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	newLogger := ctx.Logger()
	return &zerologAdapter{logger: &newLogger}
}
