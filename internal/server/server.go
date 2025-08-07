package server

import (
	"net"

	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/handler"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	pb "github.com/highonsemicolon/aura/gen/greeter"
	"github.com/highonsemicolon/aura/internal/logger"
)

func StartGRPCServer(cfg *config.Config, logAdapter logger.Logger) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logAdapter.Fatal("failed to listen", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)
	pb.RegisterGreeterServer(s, handler.NewGreeterHandler())

	logAdapter.Info("gRPC server listening on port 50051")
	if err := s.Serve(listener); err != nil {
		logAdapter.Fatal("failed to serve", err)
	}
}
