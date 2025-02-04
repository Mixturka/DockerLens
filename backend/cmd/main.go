package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Mixturka/DockerLens/backend/internal/app/config"
	"github.com/Mixturka/DockerLens/backend/internal/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log := logging.SetupLogger(config.LogLevel)
	log.Info("Starting DockerLens")
	log.Debug("Logging in Debug mode")

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// var pingRepo interfaces.PingRepository = postgres.NewPostgresPingRepository()
	fmt.Println(config)
}
