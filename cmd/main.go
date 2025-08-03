package main

import (
	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/logger"
)

func main() {

	logger := logger.New()

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Error().Msg("failed to load config: " + err.Error())
	}

	logger.Info().Msgf("Service Name: %s", cfg.ServiceName)

}
