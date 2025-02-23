package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MySQL MySQL `yaml:"mysql"`
}

type MySQL struct {
	DSN    string `yaml:"dsn"`
	CAPath string `yaml:"ca-path"`
}

var (
	instance *Config
	once     sync.Once
)

func loadConfig() {
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	instance = &config
}

func GetConfig() *Config {
	once.Do(loadConfig)
	return instance
}
