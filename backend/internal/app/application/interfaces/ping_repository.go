package interfaces

import (
	"context"

	"github.com/Mixturka/DockerLens/backend/internal/app/domain/entities"
)

type PingRepository interface {
	Save(ctx context.Context, ping entities.Ping) error
	GetById(ctx context.Context, id string) (entities.Ping, error)
	Remove(ctx context.Context, id string) error
}
