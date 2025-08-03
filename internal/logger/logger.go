package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New() *zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().
		Timestamp().
		Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return &logger
}
