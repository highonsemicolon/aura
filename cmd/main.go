package main

import (
	"context"
	"fmt"
	"net"

	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/logger"
	pb "github.com/highonsemicolon/aura/internal/proto/greeter"
	"github.com/highonsemicolon/aura/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	tracer := otel.Tracer("github.com/highonsemicolon/aura/cmd/grpc")
	_, span := tracer.Start(ctx, "SayHello")
	defer span.End()

	span.SetAttributes(
		attribute.String("method", "SayHello"),
		attribute.String("name", req.Name),
	)

	message := fmt.Sprintf("Hello, %s!", req.Name)
	return &pb.HelloResponse{Message: message}, nil
}

func main() {
	cfg := config.LoadConfig()
	logAdapter := logger.NewZerologAdapter(cfg.Logging.Format, cfg.Logging.Level)

	telemetryShutdown := telemetry.InitTracer(cfg.ServiceName, cfg.OTEL.Endpoint)
	defer telemetryShutdown()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logAdapter.Fatal("failed to listen", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)
	pb.RegisterGreeterServer(s, &server{})

	logAdapter.Info("gRPC server listening on port 50051")
	if err := s.Serve(listener); err != nil {
		logAdapter.Fatal("failed to serve", err)
	}
}
