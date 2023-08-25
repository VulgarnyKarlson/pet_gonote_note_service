package postgres

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/jackc/pgx/v4"
)

type Transaction struct {
	ctx  context.Context
	tx   pgx.Tx
	conn *pgxpool.Conn
}

func (t *Transaction) Ctx() context.Context {
	return t.ctx
}

func (t *Transaction) Exec(query string, args ...any) (pgconn.CommandTag, error) {
	return t.tx.Exec(t.ctx, query, args...)
}

func (t *Transaction) Commit() error {
	return t.tx.Commit(t.ctx)
}

func (t *Transaction) Rollback() error {
	return t.tx.Rollback(t.ctx)
}

func (t *Transaction) ConnRelease() {
	t.conn.Release()
}

func (t *Transaction) SendBatch(batch *pgx.Batch) pgx.BatchResults {
	return t.tx.SendBatch(t.ctx, batch)
}

func (t *Transaction) QueryRow(query string, args ...any) pgx.Row {
	return t.tx.QueryRow(t.ctx, query, args...)
}

func (t *Transaction) Query(query string, args ...any) (pgx.Rows, error) {
	return t.tx.Query(t.ctx, query, args...)
}
