package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Pool struct {
	pool *pgxpool.Pool
}

func New(config *Config) (*Pool, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d",
		config.UserName, config.Password, config.Host, config.Port, config.DBName, config.PoolSize)
	dbConn, err := pgxpool.Connect(context.TODO(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return &Pool{pool: dbConn}, nil
}

func (w *Pool) GetConnection() *pgxpool.Pool {
	return w.pool
}

func (w *Pool) Close() {
	w.pool.Close()
}

func (w *Pool) NewTransaction(ctx context.Context, opts pgx.TxOptions) (*Transaction, error) {
	conn, err := w.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Transaction{ctx: ctx, tx: tx, conn: conn}, nil
}
