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

	shutdownTelemetry := telemetry.InitTracer(ctx, cfg.ServiceName, cfg.OTEL.Endpoint)
	defer func() {
		_ = shutdownTelemetry(ctx)
	}()

	server.StartGRPCServer(ctx, cfg, logAdapter)
}
