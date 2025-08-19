package server

import (
	"context"
	"fmt"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	pb "github.com/highonsemicolon/aura/gen/greeter"
	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/handler"
	"github.com/highonsemicolon/aura/pkg/logging"
)

func StartGRPCServer(ctx context.Context, cfg *config.Config, log logging.Logger) error {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerZerologInterceptor(log),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	pb.RegisterGreeterServer(s, handler.NewGreeterHandler())

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("greeter", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	errCh := make(chan error, 1)
	go func() {
		log.InfoF("gRPC server listening on %s", " :50051")
		if serveErr := s.Serve(listener); serveErr != nil {
			errCh <- serveErr
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("context cancelled, shutting down gRPC server")
		healthServer.SetServingStatus("greeter", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

		s.GracefulStop()
	case serveErr := <-errCh:
		return fmt.Errorf("gRPC server error: %w", serveErr)
	}

	log.Info("gRPC server stopped gracefully")
	return nil
}
