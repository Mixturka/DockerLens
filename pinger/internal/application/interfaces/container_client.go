package interfaces

import (
	"context"
	"log/slog"
)

type ContainerClient interface {
	MonitorContainers(ctx context.Context, log *slog.Logger)
	Close()
	GetContainerIPs() ([]string, error)
}
