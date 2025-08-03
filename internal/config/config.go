package config

import (
	"strings"

	"github.com/highonsemicolon/aura/internal/logger"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ServiceName string `koanf:"service_name" validate:"required"`
	Port        int    `koanf:"port" validate:"required"`
	OTEL        struct {
		Endpoint string `koanf:"endpoint" validate:"required"`
	}
	LogLevel string `koanf:"log_level" validate:"required,oneof=debug info warn error fatal panic"`
}

var k = koanf.New(".")

func LoadConfig(logger *logger.LoggerService) (*Config, error) {

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

	logger.SetLevel(mainConfig.LogLevel)

	return mainConfig, nil
}
