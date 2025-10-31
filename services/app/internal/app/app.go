package app

import (
	"context"
	"fmt"
	"time"

	"github.com/highonsemicolon/aura/pkg/db"
	"github.com/highonsemicolon/aura/pkg/healthz"
	"github.com/highonsemicolon/aura/pkg/logging"
	"github.com/highonsemicolon/aura/pkg/telemetry"
	"github.com/highonsemicolon/aura/services/app/internal/config"
	"github.com/highonsemicolon/aura/services/app/internal/server"
)

type AppService struct {
	cfg       *config.Config
	log       logging.Logger
	healthz   *healthz.Healthz
	lifecycle *Lifecycle
}

func New(cfg *config.Config, log logging.Logger) *AppService {
	return &AppService{
		cfg:       cfg,
		log:       log,
		lifecycle: NewLifecycle(log),
	}
}

func (s *AppService) Start(ctx context.Context) error {
	shutdownTelemetry := telemetry.InitTracer(ctx, telemetry.TracerInitOption{
		ServiceName: s.cfg.ServiceName,
		Endpoint:    s.cfg.OTEL.Endpoint,
		Logger:      s.log,
	})
	s.lifecycle.Add("Telemetry", shutdownTelemetry)

	registry, mongoClient, err := db.InitMongoRegistry(ctx, s.cfg.MongoDB.URI, s.cfg.MongoDB.Databases())
	if err != nil {
		return fmt.Errorf("failed to connect db: %w", err)
	}
	s.lifecycle.Add("MongoDB Client", mongoClient.Disconnect)

	_ = registry
	// orderRepo := mongo.NewOrderRepository(registry[config.DBOrders])
	// orderRepo.CreateOrder(ctx, "order123", 99.99)

	s.healthz = healthz.NewHealthz(5 * time.Second)
	s.healthz.RegisterLiveness("liveness")
	s.healthz.RegisterReadiness("greeter", db.CheckMongoConnection(mongoClient))
	s.healthz.RegisterReadiness("thanker")
	s.healthz.Start(ctx)
	s.lifecycle.Add("Healthz", s.healthz.Stop)

	srv, err := server.New(s.cfg, s.healthz, s.log)
	if err != nil {
		return err
	}
	s.lifecycle.Add("gRPC Server", srv.Stop)

	go func() {
		if err := srv.Start(ctx); err != nil {
			s.log.Error("gRPC server error", err)
		}
	}()

	return nil
}

func (s *AppService) Shutdown(ctx context.Context) error {
	return s.lifecycle.Shutdown(ctx)
}
