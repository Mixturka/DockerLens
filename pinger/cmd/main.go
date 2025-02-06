package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Mixturka/DockerLens/backend/pkg/logging"
	"github.com/Mixturka/DockerLens/pinger/internal/application/services"
	"github.com/Mixturka/DockerLens/pinger/internal/config"
	apiclient "github.com/Mixturka/DockerLens/pinger/internal/infrastructure/api_client"
	dockerclient "github.com/Mixturka/DockerLens/pinger/internal/infrastructure/docker_client"
	iprepository "github.com/Mixturka/DockerLens/pinger/internal/infrastructure/ip_repository"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log := logging.SetupLogger(config.LogLevel)
	log.Info("Starting Pinger")
	log.Debug("Logging in Debug mode")

	apiUrl := os.Getenv("BACKEND_API_URL")
	if apiUrl == "" {
		log.Error("Backend API url was not found in environment")
		os.Exit(1)
	}
	ac := apiclient.NewApiClient(apiUrl)
	ip_repository := iprepository.NewInMemoryIpRepository()

	dc, err := dockerclient.NewDockerClient(ip_repository, config.ContainerFetchInterval)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer dc.Close()

	pinger := services.NewPingerService(dc, ac, 5*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pinger.StartPinging(ctx, log)
}
