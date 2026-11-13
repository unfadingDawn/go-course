package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MovieBaseURL     string
	HTTPAddr         string
	GrpcAddr         string
	DatabaseURL      string
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

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

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
		DatabaseURL:      databaseURL,
		UpstreamTimeout:  upstreamTimeout,
		OperationTimeout: operationTimeout,
	}

	return config, nil
}
