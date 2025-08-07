package main

import (
	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/logger"
	"github.com/highonsemicolon/aura/internal/server"
	"github.com/highonsemicolon/aura/internal/telemetry"
)

func main() {
	cfg := config.LoadConfig()
	logAdapter := logger.NewZerologAdapter(cfg.Logging.Format, cfg.Logging.Level)

	telemetryShutdown := telemetry.InitTracer(cfg.ServiceName, cfg.OTEL.Endpoint)
	defer telemetryShutdown()

	server.StartGRPCServer(cfg, logAdapter)
}
