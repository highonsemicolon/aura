package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL MySQL `yaml:"mysql" mapstructure:"mysql"`
	Tables map[string]string `yaml:"tables" mapstructure:"tables"`
}

type MySQL struct {
	DSN    string `yaml:"dsn" mapstructure:"dsn"`
	CAPath string `yaml:"ca-path" mapstructure:"ca-path"`
}

var (
	instance *Config
	once     sync.Once
)

func loadConfig() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	instance = &config
}

func GetConfig() *Config {
	once.Do(loadConfig)
	return instance
}
