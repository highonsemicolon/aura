package telemetry

import (
	"context"

	"github.com/highonsemicolon/aura/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitTracer(serviceName string) func() {
	logger := logger.NewZerologAdapter("json", "info")

	exp, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint("http://localhost:4317"),
	)
	if err != nil {
		logger.Fatal("Failed to create trace exporter", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	logger.Info("OpenTelemetry tracer initialized")

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Error("Error shutting down tracer provider", err)
		}
	}
}
