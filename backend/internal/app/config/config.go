package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddr  string
	LogLevel    string
	PostgresCfg PostgresConfig
	CorsCfg     CorsConfig
}

type PostgresConfig struct {
	User     string
	Password string
	Db       string
	Port     string
	Host     string
}

type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return strings.Split(val, ",")
}

func getEnvAsBool(key string, defaultValue bool) bool {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return boolVal
}

func getEnvAsInt(key string, defaultValue int) int {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("failed to load config from environment: %s", err)
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		return Config{}, errors.New("listen address must be explicitly set")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "local"
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
		CorsCfg: CorsConfig{
			AllowedOrigins:   getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
			AllowedMethods:   getEnvAsSlice("ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders:   getEnvAsSlice("ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type"}),
			ExposedHeaders:   getEnvAsSlice("EXPOSED_HEADERS", []string{"Link", "WWW-Authenticate"}),
			AllowCredentials: getEnvAsBool("ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvAsInt("CORS_MAX_AGE", 600),
		},
	}, nil
}
