package dockerclient

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Mixturka/DockerLens/pinger/internal/application/interfaces"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	cli     *client.Client
	ipStore interfaces.IpRepository
}

func NewDockerClient(ipStore interfaces.IpRepository, containerFetchInterval time.Duration, dockerOpts ...client.Opt) (*DockerClient, error) {
	client, err := client.NewClientWithOpts(dockerOpts...)
	if err != nil {
		return nil, err
	}
	return &DockerClient{
		cli:     client,
		ipStore: ipStore,
	}, nil
}

func (dc *DockerClient) Close() {
	dc.cli.Close()
}

func (dc *DockerClient) GetContainerIPs() ([]string, error) {
	return dc.ipStore.GetAll()
}

func (dc *DockerClient) MonitorContainers(ctx context.Context, log *slog.Logger) {
	// Required to correctly wait for self. Sometimes pinger launches before it can determine its container's IP
	selfId := os.Getenv("HOSTNAME")
	err := dc.waitForContainerRunning(log, selfId, 5*time.Second)
	if err != nil {
		log.Error("Failed to start pinger")
		return
	}

	containers, err := dc.cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return
	}

	for _, container := range containers {
		ip, err := dc.getContainerIp(container.ID)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		log.Debug("Container ip found: " + ip)
		err = dc.ipStore.Add(ip)
		if err != nil {
			log.Error(err.Error())
		}
	}
	messages, errs := dc.cli.Events(context.Background(), events.ListOptions{})

	for {
		select {
		case msg := <-messages:
			if msg.Type == events.ContainerEventType {
				switch msg.Action {
				case "start", "create":
					ip, err := dc.getContainerIp(msg.Actor.ID)
					if err != nil {
						log.Error("Failed to get started or created container IP", slog.Any("error", err.Error()))
						break
					}
					err = dc.ipStore.Add(ip)
					if err != nil {
						log.Error("Failed to add ip to ip storage", slog.Any("error", err.Error()))
					}
				case "stop", "die":
					err := dc.ipStore.Remove(msg.Actor.ID)
					if err != nil {
						log.Error("Failed to remove stopped/died container IP from ip storage", slog.Any("error", err.Error()))
					}
				}
			}
		case err := <-errs:
			if err != nil {
				log.Error("Failed to fetch container", slog.Any("error", err.Error()))
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (dc *DockerClient) getContainerIp(containerId string) (string, error) {
	containerInfo, err := dc.cli.ContainerInspect(context.Background(), containerId)
	if err != nil {
		return "", err
	}

	for _, net := range containerInfo.NetworkSettings.Networks {
		if net.IPAddress != "" {
			return net.IPAddress, nil
		}
	}

	return "", fmt.Errorf("failed to determine IP of container with ID: %s", containerId)
}

func (dc *DockerClient) waitForContainerRunning(log *slog.Logger, containerId string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			inspect, err := dc.cli.ContainerInspect(ctx, containerId)
			if err == nil && inspect.State != nil && inspect.State.Running {
				if len(inspect.NetworkSettings.Networks) > 0 {
					for _, network := range inspect.NetworkSettings.Networks {
						if network.IPAddress != "" {
							return nil
						}
					}
				}
				log.Info("Container is running but waiting for IP address")

			} else if err != nil {
				log.Error(fmt.Sprintf("Error inspecting container: %v", err))
			} else {
				log.Info("Container is not running yet")
			}

			time.Sleep(1 * time.Second)
		}
	}
}
