package config

type Config struct {
	Env         string  `koanf:"env" validate:"required,oneof=dev prod"`
	ServiceName string  `koanf:"service_name" validate:"required"`
	GRPC        GRPC    `koanf:"grpc" validate:"required"`
	OTEL        OTEL    `koanf:"otel"`
	Logging     Logging `koanf:"logging"`
	MongoDB     MongoDB `koanf:"mongodb" validate:"required"`
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
