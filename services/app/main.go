package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	c "github.com/highonsemicolon/aura/pkg/config"
	"github.com/highonsemicolon/aura/pkg/healthz"
	"github.com/highonsemicolon/aura/pkg/logging"
	"github.com/highonsemicolon/aura/pkg/telemetry"
	"github.com/highonsemicolon/aura/services/app/internal/config"
	"github.com/highonsemicolon/aura/services/app/internal/server"
)

var (
	Version   = "dev"
	Commit    = ""
	BuildTime = ""
	BuiltBy   = ""
)

func run(ctx context.Context) error {
	logAdapter := logging.NewZerologAdapter(logging.LoggingOption{
		Format: "json",
		Level:  "info",
	})

	cfg := &config.Config{}
	err := c.Load(cfg, c.ConfigLoaderOption{
		Prefix: "app.",
		Logger: logAdapter,
	})
	if err != nil {
		logAdapter.Fatal("failed to load config", err)
	}

	logAdapter = logging.NewZerologAdapter(logging.LoggingOption{
		Format: cfg.Logging.Format,
		Level:  cfg.Logging.Level,
	})

	logAdapter.Info("starting service:", cfg.ServiceName)
	logAdapter.Info("version:", Version)
	logAdapter.Info("commit:", Commit)
	logAdapter.Info("build_time:", BuildTime)
	logAdapter.Info("built_by:", BuiltBy)

	shutdownTelemetry := telemetry.InitTracer(ctx, telemetry.TracerInitOption{
		ServiceName: cfg.ServiceName,
		Endpoint:    cfg.OTEL.Endpoint,
		Logger:      logAdapter,
	})
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

	if err := server.StartGRPCServer(ctx, &cfg.GRPC, healthz, logAdapter); err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}
	return nil
}

func checkDBConnection() healthz.Checker {
	return func(ctx context.Context) bool {
		_, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		return true
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}
