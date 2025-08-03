package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type LoggerService struct {
	*zerolog.Logger
}

func New() *LoggerService {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().
		Timestamp().
		Logger()

	return &LoggerService{
		Logger: &logger,
	}
}

func (ls *LoggerService) SetLevel(level string) {
	ls.Logger.Info().Msgf("setting log level to %s", level)
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
