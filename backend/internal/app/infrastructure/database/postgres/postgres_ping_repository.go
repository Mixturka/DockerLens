package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mixturka/DockerLens/backend/internal/app/domain/entities"
	"github.com/jackc/pgx/v5"
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
	query := `INSERT INTO pings (ip, is_success, ping_time, time_stamp) 
				VALUES ($1, $2, $3, $4)`
	_, err := pr.pool.Exec(ctx, query, ping.IP, ping.IsSuccess, ping.Duration, ping.LastSuccess)

	if err != nil {
		return err
	}
	return nil
}

func (pr *PostgresPingRepository) GetAllPings(ctx context.Context) ([]entities.Ping, error) {
	query := `SELECT ip::text, is_success, ping_time, time_stamp FROM pings`
	rows, err := pr.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pings []entities.Ping
	for rows.Next() {
		var ping entities.Ping
		if err := rows.Scan(&ping.IP, &ping.IsSuccess, &ping.Duration, &ping.LastSuccess); err != nil {
			return nil, err
		}
		pings = append(pings, ping)
	}
	return pings, nil
}

func (pr *PostgresPingRepository) SaveBatch(ctx context.Context, pings []entities.Ping) error {
	if len(pings) == 0 {
		return nil
	}

	query := `INSERT INTO pings (ip, is_success, ping_time, time_stamp) VALUES `
	vals := []any{}
	valueStrings := make([]string, 0, len(pings))

	for i, ping := range pings {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		vals = append(vals, ping.IP, ping.IsSuccess, ping.Duration, ping.LastSuccess)
	}

	query += strings.Join(valueStrings, ", ")

	query += `
        ON CONFLICT (ip) DO UPDATE SET
            is_success = EXCLUDED.is_success,
            ping_time = EXCLUDED.ping_time,
            time_stamp = CASE
                WHEN EXCLUDED.is_success = true THEN EXCLUDED.time_stamp
                ELSE (SELECT time_stamp FROM pings WHERE ip = EXCLUDED.ip)
            END`

	if _, err := pr.pool.Exec(ctx, query, vals...); err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}
	return nil
}

func (pr *PostgresPingRepository) GetLatestSuccessByIp(ctx context.Context, ip string) (entities.Ping, error) {
	var ping entities.Ping
	query := `SELECT ip, is_success, ping_time, time_stamp FROM pings
			  WHERE ip = $1 AND is_success = true
			  ORDER BY time_stamp DESC LIMIT 1`

	err := pr.pool.QueryRow(ctx, query, ip).Scan(&ping.IP, &ping.IsSuccess, &ping.Duration, &ping.LastSuccess)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.Ping{}, nil
		} else {
			return entities.Ping{}, err
		}
	}

	return ping, nil
}

func (pr *PostgresPingRepository) Remove(ctx context.Context, ip string) error {
	query := `DELETE FROM pings WHERE ip = $1`
	_, err := pr.pool.Exec(ctx, query, ip)

	if err != nil {
		return err
	}
	return nil
}

func (pr *PostgresPingRepository) Healthcheck(ctx context.Context) error {
	return pr.pool.Ping(ctx)
}
