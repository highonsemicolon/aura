package main

import (
	"net/http"

	"github.com/highonsemicolon/aura/internal/config"
	"github.com/highonsemicolon/aura/internal/logger"
	"github.com/highonsemicolon/aura/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {

	cfg := config.LoadConfig()
	logger := logger.NewZerologAdapter(cfg.Logging.Format, cfg.Logging.Level)

	telemetryShutdown := telemetry.InitTracer(cfg.ServiceName)
	defer telemetryShutdown()

	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		http.ResponseWriter.WriteHeader(w, http.StatusOK)
		w.Write([]byte("Hello!"))
	}), "RootHandler"))

	logger.DebugF("service name: %s", cfg.ServiceName)

	logger.Info("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("Server error", err)
	}
}
