package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/highonsemicolon/aura/cmd/app/internal/config"
	"github.com/highonsemicolon/aura/cmd/app/internal/server"
	"github.com/highonsemicolon/aura/pkg/healthz"
	"github.com/highonsemicolon/aura/pkg/logging"
	"github.com/highonsemicolon/aura/pkg/telemetry"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	cfg := config.LoadConfig()

	logAdapter := logging.NewZerologAdapter(cfg.Logging.Format, cfg.Logging.Level)

	shutdownTelemetry := telemetry.InitTracer(ctx, cfg.ServiceName, cfg.OTEL.Endpoint)
	defer func() {
		if err := shutdownTelemetry(ctx); err != nil {
			logAdapter.Error("failed to shutdown telemetry", err)
		}
	}()

	healthz := healthz.NewHealthz(5 * time.Second)
	healthz.RegisterLiveness("liveness")
	healthz.RegisterReadiness("greeter", checkDBConnection())
	healthz.RegisterReadiness("thanker")
	healthz.Start(ctx)
	defer healthz.Stop()

	if err := server.StartGRPCServer(ctx, cfg, healthz, logAdapter); err != nil {
		logAdapter.Error("gRPC server failed", err)
		os.Exit(1)
	}
}

func checkDBConnection() healthz.Checker {
	return func(ctx context.Context) bool {
		_, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		return true
	}
}
