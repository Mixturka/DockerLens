package postgres

import (
	"context"
	"fmt"

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
	_, err := pr.pool.Exec(ctx, query, ping.ID, ping.IP, ping.IsSuccess, ping.Duration, ping.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (pr *PostgresPingRepository) SaveBatch(ctx context.Context, pings []entities.Ping) error {
	if len(pings) == 0 {
		return nil
	}
	query := `INSERT INTO pings (id, ip, is_success, ping_time, time_stamp) VALUES `
	vals := []interface{}{}

	for i, ping := range pings {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		vals = append(vals, ping.ID, ping.IP, ping.IsSuccess, ping.Duration, ping.CreatedAt)
	}

	query = query[:len(query)-1]
	_, err := pr.pool.Exec(ctx, query, vals...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}
	return nil
}

func (pr *PostgresPingRepository) GetById(ctx context.Context, id string) (entities.Ping, error) {
	var ping entities.Ping
	query := `SELECT id, ip, is_success, time_stamp from pings WHERE id = $1`

	err := pr.pool.QueryRow(ctx, query, id).Scan(&ping.ID, &ping.IP, &ping.IsSuccess, &ping.Duration, &ping.CreatedAt)

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

func (pr *PostgresPingRepository) GetPingsCursor(ctx context.Context, limit int, cursor string) ([]entities.Ping, string, error) {
	query := `SELECT id, ip, is_success, ping_time, time_stamp 
	          FROM pings `
	args := []interface{}{}

	if cursor != "" {
		query += "WHERE time_stamp < $1 "
		args = append(args, cursor)
	}

	query += "ORDER BY time_stamp DESC LIMIT $2"
	args = append(args, limit)

	rows, err := pr.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var pings []entities.Ping
	var nextCursor string

	for rows.Next() {
		var ping entities.Ping
		err := rows.Scan(&ping.ID, &ping.IP, &ping.IsSuccess, &ping.Duration, &ping.CreatedAt)
		if err != nil {
			return nil, "", err
		}
		pings = append(pings, ping)
	}

	if len(pings) > 0 {
		nextCursor = pings[len(pings)-1].CreatedAt.Format("2006-01-02T15:04:05")
	}

	return pings, nextCursor, nil
}

func (pr *PostgresPingRepository) Healthcheck(ctx context.Context) error {
	return pr.pool.Ping(ctx)
}
