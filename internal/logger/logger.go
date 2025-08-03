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
	format = strings.ToLower(format)
	level = strings.ToLower(level)

	var writer io.Writer
	if format == "json" {
		writer = os.Stderr
	} else {
		writer = zerolog.ConsoleWriter{Out: os.Stderr}
	}

	logger := zerolog.New(writer).With().
		Timestamp().
		Logger()

	ls := &LoggerService{
		Logger: &logger,
	}
	ls.setLevel(level)
	return ls
}

func (ls *LoggerService) setLevel(level string) {
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
