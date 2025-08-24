package handler

import (
	"context"
	"fmt"
	"time"

	pb "github.com/highonsemicolon/aura/apis/greeter/gen"
	"github.com/highonsemicolon/aura/pkg/logging"
	"github.com/highonsemicolon/aura/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
)

type GreeterHandler struct {
	pb.UnimplementedGreeterServer
}

func NewGreeterHandler() *GreeterHandler {
	return &GreeterHandler{}
}

func (s *GreeterHandler) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log := logging.FromContext(ctx)

	tracer := telemetry.Tracer("github.com/highonsemicolon/aura/cmd/grpc")
	_, span := tracer.Start(ctx, "SayHello")
	defer span.End()

	time.Sleep(5000 * time.Millisecond)

	span.SetAttributes(
		attribute.String("method", "SayHello"),
		attribute.String("name", req.Name),
	)

	log.InfoF("Received SayHello request for %s", req.Name)

	message := fmt.Sprintf("Hello, %s!", req.Name)
	return &pb.HelloResponse{Message: message}, nil
}
