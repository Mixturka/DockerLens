package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Mixturka/DockerLens/backend/internal/app/application/interfaces"
	"github.com/Mixturka/DockerLens/backend/internal/app/config"
	"github.com/Mixturka/DockerLens/backend/internal/app/infrastructure/database/postgres"
	server "github.com/Mixturka/DockerLens/backend/internal/app/infrastructure/httpserver"
	"github.com/Mixturka/DockerLens/backend/internal/app/infrastructure/httpserver/handlers/url/get"
	"github.com/Mixturka/DockerLens/backend/internal/app/infrastructure/httpserver/handlers/url/save"
	"github.com/Mixturka/DockerLens/backend/pkg/dbutils"
	"github.com/Mixturka/DockerLens/backend/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log := logging.SetupLogger(config.LogLevel)
	log.Info("Starting DockerLens backend")
	log.Debug("Logging in Debug mode")

	connStr := dbutils.BuildPostgresURL(config.PostgresCfg.User,
		config.PostgresCfg.Password,
		config.PostgresCfg.Host,
		config.PostgresCfg.Port, config.PostgresCfg.Db)

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Error("Unable to create connection pool, ", slog.Any("error:", err.Error()))
		os.Exit(1)
	}
	defer dbpool.Close()

	var pingRepo interfaces.PingRepository = postgres.NewPostgresPingRepository(dbpool)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := pingRepo.Healthcheck(ctx); err != nil {
		log.Error("DB healthcheck failed: " + err.Error())
		os.Exit(1)
	}

	var router chi.Router = chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.CorsCfg.AllowedOrigins,
		AllowedMethods:   config.CorsCfg.AllowedMethods,
		AllowedHeaders:   config.CorsCfg.AllowedHeaders,
		ExposedHeaders:   config.CorsCfg.ExposedHeaders,
		AllowCredentials: config.CorsCfg.AllowCredentials,
		MaxAge:           config.CorsCfg.MaxAge,
	}))

	router.Post("/api/v1/pings", save.NewSaveHandler(log, pingRepo))
	router.Get("/api/v1/pings", get.NewGetHandler(log, pingRepo))

	server := server.NewServer(router)
	if err := server.Start(log, config); err != nil {
		log.Error("Failed to start server")
		os.Exit(1)
	}
	defer server.Shutdown(log)
}
