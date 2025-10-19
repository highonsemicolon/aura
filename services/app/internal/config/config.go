package config

type Config struct {
	ServiceName string  `koanf:"service_name" validate:"required"`
	GRPC        GRPC    `koanf:"grpc" validate:"required"`
	OTEL        OTEL    `koanf:"otel"`
	Logging     Logging `koanf:"logging"`
}

type GRPC struct {
	Address string `koanf:"address" validate:"required"`
}

type Logging struct {
	Level  string `koanf:"level" validate:"required,oneof=debug info warn error fatal panic"`
	Format string `koanf:"format" validate:"required,oneof=json console"`
}

type OTEL struct {
	Endpoint string `koanf:"endpoint" validate:"required"`
}
