package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/handler"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	pb "github.com/highonsemicolon/aura/gen/greeter"
	"github.com/highonsemicolon/aura/internal/logger"
)

func StartGRPCServer(ctx context.Context, cfg *config.Config, logger logger.Logger) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal("failed to listen", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)
	pb.RegisterGreeterServer(s, handler.NewGreeterHandler())

	go func() {
		logger.Info("gRPC server listening on port 50051")
		if err := s.Serve(listener); err != nil {
			logger.Fatal("failed to serve", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logger.Info("context cancelled, shutting down")
	case sig := <-stop:
		logger.InfoF("received signal: %s, shutting down", sig)
	}

	s.GracefulStop()
	logger.Info("gRPC server stopped gracefully")
}
