package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func New(config *Config) (*pgxpool.Pool, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d",
		config.UserName, config.Password, config.Host, config.Port, config.DBName, config.PoolSize)
	dbConn, err := pgxpool.Connect(context.TODO(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return dbConn, nil
}
