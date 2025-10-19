package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var tracerProvider *sdktrace.TracerProvider

type TracerInitOption struct {
	ServiceName string
	Endpoint    string
	Logger      Logger
}

type Logger interface {
	Info(msg ...string)
	Fatal(msg string, errs ...error)
}

func InitTracer(ctx context.Context, opts TracerInitOption) func(context.Context) error {

	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(opts.Endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		opts.Logger.Fatal("failed to create trace exporter", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(opts.ServiceName),
		)),
	)

	tracerProvider = tp
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	opts.Logger.Info("open telemetry tracer initialized")

	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			return err
		}
		return exporter.Shutdown(ctx)
	}
}

func Tracer(name string) trace.Tracer {
	return tracerProvider.Tracer(name)
}
