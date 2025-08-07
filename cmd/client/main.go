package main

import (
	"context"
	"fmt"
	"os"

	pb "github.com/highonsemicolon/aura/gen/greeter"
	"github.com/highonsemicolon/aura/internal/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318"
	}

	telemetryShutdown := telemetry.InitTracer("grpc-client-service", endpoint)
	defer telemetryShutdown()

	tracer := otel.Tracer("grpc-client")

	ctx, span := tracer.Start(context.Background(), "call-grpc-server")
	defer span.End()

	span.SetAttributes(
		attribute.String("client", "my-grpc-client"),
	)

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("connection_error", err.Error()))
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Aura"})
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("grpc_error", err.Error()))
		fmt.Println("gRPC request failed:", err)
		return
	}

	fmt.Println("Response:", resp.Message)
}
