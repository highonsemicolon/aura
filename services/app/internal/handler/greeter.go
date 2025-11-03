package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/highonsemicolon/aura/apis/gen/greeter"
	"github.com/highonsemicolon/aura/pkg/logging"
	"github.com/highonsemicolon/aura/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
)

type GreeterHandler struct {
	greeter.UnimplementedGreeterServiceServer
}

func NewGreeterHandler() *GreeterHandler {
	return &GreeterHandler{}
}

func (s *GreeterHandler) SayHello(ctx context.Context, req *greeter.SayHelloRequest) (*greeter.HelloResponse, error) {
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
	return &greeter.HelloResponse{Message: message}, nil
}
