package config

import (
	"strings"

	"github.com/go-playground/validator"
	"github.com/highonsemicolon/aura/internal/logger"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ServiceName string  `koanf:"service_name" validate:"required"`
	Port        int     `koanf:"port" validate:"required"`
	OTEL        OTEL    `koanf:"otel"`
	Logging     Logging `koanf:"logging"`
}

type Logging struct {
	Level  string `koanf:"level" validate:"required,oneof=debug info warn error fatal panic"`
	Format string `koanf:"format" validate:"required,oneof=json console"`
}

type OTEL struct {
	Endpoint string `koanf:"endpoint" validate:"required"`
}

var k = koanf.New(".")

func LoadConfig() *Config {
	logger := logger.NewZerologAdapter("json", "info")

	err := k.Load(env.Provider("TEMPLATE_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "TEMPLATE_"))
	}), nil)
	if err != nil {
		logger.Fatal("could not load initial env variables")
	}

	mainConfig := &Config{}

	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal("could not unmarshal main config", err)
	}

	validate := validator.New()

	err = validate.Struct(mainConfig)
	if err != nil {
		logger.Fatal("config validation failed", err)
	}

	return mainConfig
}
