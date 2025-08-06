package main

import (
	"net/http"

	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/logger"
	"github.com/highonsemicolon/aura/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
)

func main() {

	cfg := config.LoadConfig()
	logger := logger.NewZerologAdapter(cfg.Logging.Format, cfg.Logging.Level)

	telemetryShutdown := telemetry.InitTracer(cfg.ServiceName, cfg.OTEL.Endpoint)
	defer telemetryShutdown()

	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tracer := telemetry.Tracer("github.com/highonsemicolon/aura/cmd/main")
		_, span := tracer.Start(r.Context(), "processing-root-request")
		defer span.End()

		span.SetAttributes(
			attribute.String("handler", "root"),
			attribute.String("method", r.Method),
		)

		http.ResponseWriter.WriteHeader(w, http.StatusOK)
		w.Write([]byte("Hello!"))
	}), "RootHandler"))

	logger.DebugF("service name: %s", cfg.ServiceName)

	logger.Info("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("Server error", err)
	}
}
