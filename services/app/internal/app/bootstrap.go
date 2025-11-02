package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/highonsemicolon/aura/pkg/configloader"
	"github.com/highonsemicolon/aura/pkg/logging"
	"github.com/highonsemicolon/aura/services/app/internal/config"
)

func RunApp(version, commit, buildTime, builtBy string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logging.NewZerologAdapter(logging.LoggingOption{
		Format: "json",
		Level:  "info",
	})

	log.InfoF("version=%s commit=%s buildTime=%s builtBy=%s", version, commit, buildTime, builtBy)

	defer func() {
		if r := recover(); r != nil {
			log.FatalF("unhandled panic: %v\n%s", r, debug.Stack())
		}
	}()

	// Load configuration
	cfg := &config.Config{}
	if err := configloader.Load(cfg, configloader.ConfigLoaderOption{
		Prefix: "app.",
		Logger: log,
	}); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	svc := New(cfg, log)
	if err := svc.Start(ctx); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	<-ctx.Done()
	log.Info("shutting down service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := svc.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	} else {
		log.Info("service stopped gracefully")
	}

	return nil
}
