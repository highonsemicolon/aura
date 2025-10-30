package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	c "github.com/highonsemicolon/aura/pkg/config"
	"github.com/highonsemicolon/aura/pkg/logging"
	"github.com/highonsemicolon/aura/services/app/internal/config"
)

var (
	Version   = "dev"
	Commit    = ""
	BuildTime = ""
	BuiltBy   = ""
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logging.NewZerologAdapter(logging.LoggingOption{
		Format: "json",
		Level:  "info",
	})

	cfg := &config.Config{}
	err := c.Load(cfg, c.ConfigLoaderOption{
		Prefix: "app.",
		Logger: log,
	})
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	service := NewAppService(cfg, log)
	if err := service.Start(ctx); err != nil {
		log.Fatal("failed to start service", err)
	}

	<-ctx.Done()
	log.Info("shutting down service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", err)
	} else {
		log.Info("service stopped gracefully")
	}
}
