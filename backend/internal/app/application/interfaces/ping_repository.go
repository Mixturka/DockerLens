package interfaces

import (
	"context"

	"github.com/Mixturka/DockerLens/backend/internal/app/domain/entities"
)

type PingRepository interface {
	Save(ctx context.Context, ping entities.Ping) error
	SaveBatch(ctx context.Context, pings []entities.Ping) error
	GetAllPings(ctx context.Context) ([]entities.Ping, error)
	GetLatestSuccessByIp(ctx context.Context, ip string) (entities.Ping, error)
	Remove(ctx context.Context, iз string) error
	Healthcheck(ctx context.Context) error
}
