package services

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Mixturka/DockerLens/pinger/internal/application/interfaces"
	"github.com/Mixturka/DockerLens/pinger/internal/domain/entities"
	"github.com/go-ping/ping"
)

type PingerService struct {
	containerClient interfaces.ContainerClient
	apiClient       interfaces.ApiClient
	pingInterval    time.Duration
}

func NewPingerService(containerClient interfaces.ContainerClient, apiClient interfaces.ApiClient, pingInterval time.Duration) *PingerService {
	return &PingerService{
		containerClient: containerClient,
		apiClient:       apiClient,
		pingInterval:    pingInterval,
	}
}

func (ps *PingerService) StartPinging(ctx context.Context, log *slog.Logger) {
	log.Info("Starting to ping containers")
	go ps.containerClient.MonitorContainers(ctx, log)

	ticker := time.NewTicker(ps.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Info("Saving current pings")
			pingResults := ps.PingAll(log)
			err := ps.SavePingResults(pingResults, log)
			if err != nil {
				log.Error(fmt.Sprintf("Saving failed: %s", err))
			}
		}
	}
}

func (ps *PingerService) ping(log *slog.Logger, ip string) (time.Duration, error) {
	start := time.Now()

	pinger, err := ping.NewPinger(ip)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to create ICMP pinger for IP %s: %v", ip, err))
		return -1, err
	}
	pinger.Count = 1
	pinger.Timeout = 2 * time.Second

	err = pinger.Run()
	if err != nil {
		log.Error(fmt.Sprintf("Error running ping to %s: %v", ip, err))
		return -1, err
	}

	pingDuration := time.Since(start)

	return pingDuration, nil
}

func (ps *PingerService) PingAll(log *slog.Logger) []entities.PingResult {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []entities.PingResult
	time.Sleep(10 * time.Second)

	ips, err := ps.containerClient.GetContainerIPs()
	if err != nil {
		log.Error("Failed to get container IPs: ", slog.Any("error", err))
		return []entities.PingResult{}
	}
	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			duration, err := ps.ping(log, ip)
			success := true
			if err != nil {
				log.Error(fmt.Sprintf("Ping failed for IP %s", ip))
				success = false
			}
			if success {
				log.Debug(fmt.Sprintf("Successfully pinged container with IP: %s", ip))
			}

			mu.Lock()
			results = append(results, entities.PingResult{
				IP:        ip,
				IsSuccess: success,
				Duration:  duration.Milliseconds(),
			})
			mu.Unlock()
		}(ip)
	}

	wg.Wait()
	return results
}

func (ps *PingerService) SavePingResults(pingResults []entities.PingResult, log *slog.Logger) error {
	return ps.apiClient.SavePingResults(pingResults)
}
