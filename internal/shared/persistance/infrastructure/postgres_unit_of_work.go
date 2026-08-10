package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresUnitOfWork is the pgx-backed UnitOfWork: it opens a transaction, puts
// it in the context (so repositories join it via TxFromContext) and commits or
// rolls back around fn. It is tx-aware — if a transaction is already in the
// context it joins it and lets the owner commit/roll back, avoiding nesting.
type PostgresUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewPostgresUnitOfWork(pool *pgxpool.Pool) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{pool: pool}
}

func (u *PostgresUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := TxFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op after commit; also covers panics

	if err = fn(ContextWithTx(ctx, tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
