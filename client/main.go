package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func initTracer(serviceName, endpoint string) (func(), error) {
	ctx := context.Background()

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() {
		_ = tp.Shutdown(ctx)
	}, nil
}

func main() {
	// Config
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318"
	}

	shutdown, err := initTracer("http-client-service", endpoint)
	if err != nil {
		panic(err)
	}
	defer shutdown()

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   5 * time.Second,
	}

	tracer := otel.Tracer("http-client")

	ctx, span := tracer.Start(context.Background(), "call-traced-server")
	defer span.End()

	span.SetAttributes(
		attribute.String("client", "my-http-client"),
	)

	// Make the request to the server service
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/", nil)

	res, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", err.Error()))
		fmt.Println("Request failed:", err)
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Printf("Response: %s\n", string(body))
}
