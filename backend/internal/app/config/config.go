package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Mixturka/vm-hub/pkg/putils"
	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddr  string
	LogLevel    string
	PostgresCfg PostgresConfig
}

type PostgresConfig struct {
	User     string
	Password string
	Db       string
	Port     string
	Host     string
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("failed to load config from .env: %s", err)
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		return Config{}, errors.New("listen address must be explicitly set")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "local"
	}
	cwd, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("failed to determine current working directory: %s", err)
	}
	projRoot, err := putils.GetProjectRoot(cwd)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load postgres config from database.env: %s", err)
	}
	if err := godotenv.Load(filepath.Join(projRoot, "deployments/env/database.env")); err != nil {
		return Config{}, fmt.Errorf("failed to load postgres config from database.env: %s", err)
	}

	return Config{
		ListenAddr: listenAddr,
		LogLevel:   logLevel,
		PostgresCfg: PostgresConfig{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Db:       os.Getenv("POSTGRES_DB"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Host:     os.Getenv("POSTGRES_HOST"),
		},
	}, nil
}
