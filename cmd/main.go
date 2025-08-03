package main

import (
	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/logger"
)

func main() {

	cfg := config.LoadConfig()
	logger := logger.New(cfg.Logging.Format, cfg.Logging.Level)

	logger.Debug().Msgf("service name: %s", cfg.ServiceName)
}
