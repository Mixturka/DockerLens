package interfaces

import (
	"context"
	"log/slog"

	"github.com/Mixturka/DockerLens/pinger/internal/domain/entities"
)

type ContainerPinger interface {
	PingAll() []entities.PingResult
	SavePingResults(pingResults []entities.PingResult) error
	StartPinging(ctx context.Context, log *slog.Logger)
}
