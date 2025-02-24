package config

import (
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address string `envconfig:"ADDRESS"`
	MySQL   MySQL
	Tables  map[string]string `envconfig:"TABLES"`
}

type MySQL struct {
	DSN             string        `envconfig:"MYSQL_DSN"`
	CAPath          string        `envconfig:"MYSQL_CA_PATH"`
	MaxOpenConns    int           `envconfig:"MYSQL_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `envconfig:"MYSQL_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `envconfig:"MYSQL_CONN_MAX_LIFETIME"`
}

var (
	instance *Config
	once     sync.Once
)

func loadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	instance = &config
}

func GetConfig() *Config {
	once.Do(loadConfig)
	return instance
}
