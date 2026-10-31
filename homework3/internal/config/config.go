package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MovieBaseURL     string
	HTTPAddr         string
	GrpcAddr         string
	UpstreamTimeout  time.Duration
	OperationTimeout time.Duration
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	upstreamTimeout, err := time.ParseDuration(envOrDefault("UPSTREAM_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid UPSTREAM_TIMEOUT: %w", err)
	}
	if upstreamTimeout <= 0 {
		return nil, fmt.Errorf("UPSTREAM_TIMEOUT must be positive")
	}

	operationTimeout, err := time.ParseDuration(envOrDefault("OPERATION_TIMEOUT", "30s"))
	if err != nil {
		return nil, fmt.Errorf("invalid OPERATION_TIMEOUT: %w", err)
	}
	if operationTimeout <= 0 {
		return nil, fmt.Errorf("OPERATION_TIMEOUT must be positive")
	}

	config := &Config{
		MovieBaseURL:     envOrDefault("MOVIE_BASE_URL", "https://homeworksite.site"),
		HTTPAddr:         envOrDefault("HTTP_ADDR", ":8080"),
		GrpcAddr:         envOrDefault("GRPC_ADDR", ":9090"),
		UpstreamTimeout:  upstreamTimeout,
		OperationTimeout: operationTimeout,
	}

	return config, nil
}
