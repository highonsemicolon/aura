package server

import (
	"context"
	"fmt"
	"net"

	pb "github.com/highonsemicolon/aura/apis/greeter/gen"

	"github.com/highonsemicolon/aura/pkg/healthz"
	"github.com/highonsemicolon/aura/pkg/logging"
	"github.com/highonsemicolon/aura/services/app/internal/config"
	"github.com/highonsemicolon/aura/services/app/internal/handler"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	cfg      *config.Config
	grpcSrv  *grpc.Server
	listener net.Listener
	log      logging.Logger
	healthz  *healthz.Healthz
}

func New(cfg *config.Config, healthz *healthz.Healthz, log logging.Logger) (*Server, error) {
	listener, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", cfg.GRPC.Address, err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerZerologInterceptor(log),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	pb.RegisterGreeterServer(s, handler.NewGreeterHandler())

	grpc_health_v1.RegisterHealthServer(s, healthz.Server())

	return &Server{
		cfg:      cfg,
		grpcSrv:  s,
		listener: listener,
		log:      log,
		healthz:  healthz,
	}, nil
}

func (srv *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		srv.log.InfoF("gRPC server listening on %s", srv.cfg.GRPC.Address)
		if err := srv.grpcSrv.Serve(srv.listener); err != nil {
			errCh <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		srv.log.Info("context cancelled, initiating graceful gRPC shutdown...")
		srv.healthz.SetAllNotServing()
		srv.grpcSrv.GracefulStop()
		return nil
	case err := <-errCh:
		return err
	}
}

func (srv *Server) Stop(ctx context.Context) error {
	srv.healthz.SetAllNotServing()

	done := make(chan struct{})
	go func() {
		srv.grpcSrv.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		srv.log.Info("gRPC server stopped gracefully")
		return nil
	case <-ctx.Done():
		srv.log.Warn("gRPC shutdown timed out; forcing stop")
		srv.grpcSrv.Stop()
		return ctx.Err()
	}
}
