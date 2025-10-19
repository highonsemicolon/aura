package config

import (
	"strings"

	"github.com/go-playground/validator"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

type ConfigLoaderOption struct {
	Prefix string
	Logger Logger
}

type Logger interface {
	Info(msg ...string)
	Fatal(msg string, errs ...error)
}

func Load(out any, opts ConfigLoaderOption) error {
	k := koanf.New(".")
	prefix := opts.Prefix
	if prefix == "" {
		prefix = "APP."
	}

	err := k.Load(env.Provider(prefix, ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, prefix))
	}), nil)
	if err != nil {
		opts.Logger.Fatal("could not load env variables", err)
	}

	if err := k.Unmarshal("", out); err != nil {
		opts.Logger.Fatal("could not unmarshal config", err)
	}

	validate := validator.New()
	if err := validate.Struct(out); err != nil {
		opts.Logger.Fatal("config validation failed", err)
	}

	opts.Logger.Info("config loaded successfully")
	return nil
}
