package main

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

	mongox "go.mongodb.org/mongo-driver/mongo"
)

type AppService struct {
	cfg         *config.Config
	log         logging.Logger
	healthz     *healthz.Healthz
	shutdownFns []ShutdownFn
}

type ShutdownFn struct {
	Name string
	Fn   func(ctx context.Context) error
}

func NewAppService(cfg *config.Config, log logging.Logger) *AppService {
	return &AppService{cfg: cfg, log: log}
}

func (s *AppService) Start(ctx context.Context) error {

	s.log.Info("starting service:", s.cfg.ServiceName)
	s.log.Info("version:", Version)
	s.log.Info("commit:", Commit)
	s.log.Info("build_time:", BuildTime)
	s.log.Info("built_by:", BuiltBy)

	shutdownTelemetry := telemetry.InitTracer(ctx, telemetry.TracerInitOption{
		ServiceName: s.cfg.ServiceName,
		Endpoint:    s.cfg.OTEL.Endpoint,
		Logger:      s.log,
	})
	s.shutdownFns = append(s.shutdownFns, ShutdownFn{
		Name: "Telemetry",
		Fn: func(ctx context.Context) error {
			return shutdownTelemetry(ctx)
		},
	})

	registry, mongoClient, err := db.InitMongoRegistry(ctx, s.cfg.MongoDB.URI, s.cfg.MongoDB.Databases())
	if err != nil {
		return fmt.Errorf("failed to connect db: %w", err)
	}
	s.shutdownFns = append(s.shutdownFns, ShutdownFn{
		Name: "MongoDB Client",
		Fn: func(ctx context.Context) error {
			return mongoClient.Disconnect(ctx)
		},
	})

	s.healthz = healthz.NewHealthz(5 * time.Second)
	s.healthz.RegisterLiveness("liveness")
	s.healthz.RegisterReadiness("greeter", checkDBConnection(mongoClient))
	s.healthz.RegisterReadiness("thanker")
	s.healthz.Start(ctx)

	_ = registry
	// order_repo := mongo.NewOrderRepository(registry[config.DBOrders])
	// order_repo.CreateOrder(ctx, "order123", 99.99)

	srv, err := server.New(s.cfg, s.healthz, s.log)
	if err != nil {
		return err
	}

	s.shutdownFns = append(s.shutdownFns,
		ShutdownFn{
			Name: "gRPC Server",
			Fn: func(ctx context.Context) error {
				srv.Stop()
				return nil
			},
		})

	go func() {
		if err := srv.Start(ctx); err != nil {
			s.log.Error("gRPC server exited with error", err)
		}
	}()

	return nil
}

func checkDBConnection(mongoClient *mongox.Client) healthz.Checker {
	return func(ctx context.Context) bool {
		_, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		if err := mongoClient.Ping(ctx, nil); err != nil {
			return false
		}

		return true
	}
}

func (svc *AppService) Shutdown(ctx context.Context) error {
	svc.log.Info("stopping health checks")
	svc.healthz.Stop()

	// run shutdowns in LIFO order
	for i := len(svc.shutdownFns) - 1; i >= 0; i-- {
		step := svc.shutdownFns[i]
		svc.log.InfoF("%s: shutting down", step.Name)
		if err := step.Fn(ctx); err != nil {
			svc.log.ErrorF("error while shutting down %s: %v", step.Name, err)
		} else {
			svc.log.InfoF("%s: shut down successfully", step.Name)
		}
	}

	return nil
}
