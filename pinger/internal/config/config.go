package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel               string
	ContainerFetchInterval time.Duration
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("failed to load config from environment: %s", err)
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "local"
	}

	containerFetchIntervalStr := os.Getenv("CONTAINER_FETCH_INTERVAL_MS")
	if containerFetchIntervalStr == "" {
		containerFetchIntervalStr = "3000ms"
	}
	containerFetchInterval, err := time.ParseDuration(containerFetchIntervalStr)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse container interval from environment: %s", err)
	}
	return Config{
		LogLevel:               logLevel,
		ContainerFetchInterval: containerFetchInterval,
	}, nil
}
