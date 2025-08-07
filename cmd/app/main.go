package main

import (
	"context"

	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/logger"
	"github.com/highonsemicolon/aura/internal/server"
	"github.com/highonsemicolon/aura/internal/telemetry"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.LoadConfig()
	logAdapter := logger.NewZerologAdapter(cfg.Logging.Format, cfg.Logging.Level)

	telemetryShutdown := telemetry.InitTracer(cfg.ServiceName, cfg.OTEL.Endpoint)
	defer telemetryShutdown()

	server.StartGRPCServer(ctx, cfg, logAdapter)
}
