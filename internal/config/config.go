package config

import (
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"

	_ "github.com/joho/godotenv/autoload"

	"github.com/rs/zerolog"
)

type Config struct {
	ServiceName string `koanf:"service_name" validate:"required"`
	Port        int    `koanf:"port" validate:"required"`
	OTEL        struct {
		Endpoint string `koanf:"endpoint" validate:"required"`
	}
}

var k = koanf.New(".")

func LoadConfig() (*Config, error) {

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	err := k.Load(env.Provider("BOILERPLATE_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "BOILERPLATE_"))
	}), nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}

	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal main config")
	}

	return mainConfig, nil
}
