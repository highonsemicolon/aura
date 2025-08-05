package main

import (
	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/logger"
)

func main() {

	cfg := config.LoadConfig()
	logger := logger.NewZerologAdapter(cfg.Logging.Format, cfg.Logging.Level)

	logger.DebugF("service name: %s", cfg.ServiceName)
}
