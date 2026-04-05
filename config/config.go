package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/juanMaAV92/go-utils/env"
	"github.com/juanMaAV92/go-utils/telemetry"
)

// ServiceName is the canonical name of this service.
// Used for logging, tracing, and the API route prefix.
const ServiceName = "go-echo-blueprint"

// Config holds all configuration for the application.
// Populated once at startup from environment variables.
type Config struct {
	ServiceName  string
	Environment  string
	Port         string
	GracefulTime time.Duration
	Telemetry    telemetry.Config
}

// Load reads configuration from environment variables.
// In local environment it also loads a .env file if present.
func Load() Config {
	environment := env.GetEnvironment()

	if environment == env.LocalEnvironment {
		if err := godotenv.Load(); err != nil {
			log.Println("no .env file found, continuing with environment variables")
		}
	}

	return Config{
		ServiceName:  ServiceName,
		Environment:  environment,
		Port:         env.GetEnvWithDefault("PORT", "8080"),
		GracefulTime: env.GetEnvAsDurationWithDefault("GRACEFUL_TIME", 10*time.Second),
		Telemetry: telemetry.Config{
			ServiceName: ServiceName,
			Environment: environment,
			Endpoint:    env.GetEnv("OTEL_EXPORTER_ENDPOINT"),
			Insecure:    env.GetEnvAsBoolWithDefault("OTEL_INSECURE", true),
			SampleRate:  1.0,
		},
	}
}
