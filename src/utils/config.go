package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MySQL struct {
		DSN        string `yaml:"dsn"`
		CACertPath string `yaml:"ca_cert_path"`
	} `yaml:"mysql"`
}

func LoadConfig(filename string) *Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("failed to read config file '%s': %w", filename, err))
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config data: %w", err))
	}

	return &config
}
