module github.com/highonsemicolon/aura/services/app

go 1.25.0

require (
	github.com/highonsemicolon/aura/apis/greeter v0.0.0-20250831154025-5cdd2eae1524
	github.com/highonsemicolon/aura/pkg/config v0.0.0-20251023045624-1468ddf11295
	github.com/highonsemicolon/aura/pkg/db v0.0.0-20251028181753-5322e99ff718
	github.com/highonsemicolon/aura/pkg/healthz v0.0.0-20250831154025-5cdd2eae1524
	github.com/highonsemicolon/aura/pkg/logging v0.0.0-20250831154025-5cdd2eae1524
	github.com/highonsemicolon/aura/pkg/telemetry v0.0.0-20250831154025-5cdd2eae1524
	go.mongodb.org/mongo-driver v1.17.6
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0
	go.opentelemetry.io/otel v1.38.0
	google.golang.org/grpc v1.76.0
)

require (
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.1 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250825161204-c5933d9347a5 // indirect
)

require (
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator v9.31.0+incompatible // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250825161204-c5933d9347a5 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20250818200422-3122310a409c
