package main

import (
	"context"
	"os"

	"github.com/highonsemicolon/aura/apis/gen/greeter"
	"github.com/highonsemicolon/aura/pkg/logging"
	"github.com/highonsemicolon/aura/pkg/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"go.opentelemetry.io/otel/attribute"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318"
	}
	ctx := context.Background()
	log := logging.FromContext(ctx)

	shutdownTelemetry := telemetry.InitTracer(ctx, telemetry.TracerInitOption{
		Endpoint:    endpoint,
		ServiceName: "grpc-client-service",
		Logger:      log,
	})
	defer func() {
		_ = shutdownTelemetry(ctx)
	}()

	tracer := telemetry.Tracer("grpc-client")

	ctx, span := tracer.Start(ctx, "call-grpc-server")
	defer span.End()

	span.SetAttributes(
		attribute.String("client", "my-grpc-client"),
	)

	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("connection_error", err.Error()))
		log.Error("failed to connect:", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			span.RecordError(err)
			log.Error("failed to close gRPC connection:", err)
		}
	}()

	client := greeter.NewGreeterClient(conn)

	resp, err := client.SayHello(ctx, &greeter.HelloRequest{Name: "Aura"})
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("grpc_error", err.Error()))
		log.Error("gRPC request failed:", err)
		return
	}

	log.Info("gRPC response received:", resp.Message)
}
