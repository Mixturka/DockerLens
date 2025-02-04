package postgres

import (
	"context"

	"github.com/Mixturka/DockerLens/backend/internal/app/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPingRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPingRepository(pool *pgxpool.Pool) *PostgresPingRepository {
	return &PostgresPingRepository{
		pool: pool,
	}
}

func (pr *PostgresPingRepository) Save(ctx context.Context, ping entities.Ping) error {
	query := `INSERT INTO pings (id, ip, is_success, ping_time, time_stamp) 
				VALUES ($1, $2, $3, $4, $5)`
	_, err := pr.pool.Exec(ctx, query, ping.ID, ping.IP, ping.IsSuccess, ping.Time)

	if err != nil {
		return err
	}
	return nil
}

func (pr *PostgresPingRepository) GetById(ctx context.Context, id string) (entities.Ping, error) {
	var ping entities.Ping
	query := `SELECT id, ip, is_success, time_stamp from pings WHERE id = $1`

	err := pr.pool.QueryRow(ctx, query, id).Scan(&ping.ID, &ping.IP, &ping.IsSuccess, &ping.Time, &ping.CreatedAt)

	if err != nil {
		return entities.Ping{}, err
	}
	return ping, nil
}

func (pr *PostgresPingRepository) Remove(ctx context.Context, id string) error {
	query := `DELETE FROM pings WHERE id = $1`
	_, err := pr.pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}
	return nil
}
