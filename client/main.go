package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/highonsemicolon/aura/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318"
	}

	telemetryShutdown := telemetry.InitTracer("http-client-service", endpoint)
	defer telemetryShutdown()

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
