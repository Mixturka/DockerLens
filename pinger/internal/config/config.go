package config

import (
	"fmt"
	"os"
	"strconv"
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

	containerFetchIntervalStr := os.Getenv("CONTAINER_FETCH_INTERVAL")
	if containerFetchIntervalStr == "" {
		containerFetchIntervalStr = "30000"
	}
	containerFetchInterval, err := strconv.ParseInt(containerFetchIntervalStr, 10, 64)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse container interval from environment: %s", err)
	}
	return Config{
		LogLevel:               logLevel,
		ContainerFetchInterval: time.Duration(containerFetchInterval),
	}, nil
}
