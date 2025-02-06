package interfaces

import (
	"context"

	"github.com/Mixturka/DockerLens/backend/internal/app/domain/entities"
)

type PingRepository interface {
	Save(ctx context.Context, ping entities.Ping) error
	SaveBatch(ctx context.Context, pings []entities.Ping) error
	GetById(ctx context.Context, id string) (entities.Ping, error)
	GetPingsCursor(ctx context.Context, limit int, cursor string) ([]entities.Ping, string, error)
	Remove(ctx context.Context, id string) error
	Healthcheck(ctx context.Context) error
}
