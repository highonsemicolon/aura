package app

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/highonsemicolon/aura/pkg/configloader"
	"github.com/highonsemicolon/aura/pkg/logging"
	"github.com/highonsemicolon/aura/services/app/internal/config"
)

func RunApp(version, commit, buildTime, builtBy string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logging.NewZerologAdapter(logging.LoggingOption{
		Format: "json",
		Level:  "info",
	})

	log.Info("version:", version)
	log.Info("commit:", commit)
	log.Info("build_time:", buildTime)
	log.Info("built_by:", builtBy)

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
		log.Fatal("failed to load config", err)
	}

	svc := New(cfg, log)
	if err := svc.Start(ctx); err != nil {
		log.Fatal("failed to start service", err)
	}

	<-ctx.Done()
	log.Info("shutting down service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := svc.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", err)
	} else {
		log.Info("service stopped gracefully")
	}
}
