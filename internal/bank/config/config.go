package config

import (
	"os"
)

type Config struct {
	DatabaseURL  string
	JWTSecret    string
	ServiceName  string
	OTelEndpoint string
	Port         string

	TemporalHost      string
	TemporalPort      string
	TemporalTaskQueue string
}

var Values Config

func Init() {
	Values = Config{
		DatabaseURL:  envOrDefault("DATABASE_URL", "postgres://gotrainer:verysecret@localhost:5432/gobank?sslmode=disable"),
		JWTSecret:    envOrDefault("JWT_SECRET", "super-secret-for-training-only"),
		ServiceName:  envOrDefault("SERVICE_NAME", "bank-server"),
		OTelEndpoint: envOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318"),
		Port:         envOrDefault("PORT", "8082"),

		TemporalHost:      envOrDefault("TEMPORAL_HOST", "localhost"),
		TemporalPort:      envOrDefault("TEMPORAL_PORT", "7233"),
		TemporalTaskQueue: envOrDefault("TEMPORAL_TASK_QUEUE", "bank-transfer-queue"),
	}
}

func envOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
