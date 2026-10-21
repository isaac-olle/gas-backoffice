package infrastructure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (r *PostgresRepository[T]) querier(ctx context.Context) Querier {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return r.db
}

type PostgresRepository[T any] struct {
	db *pgxpool.Pool
}

func NewPostgresRepository[T any](db *pgxpool.Pool) *PostgresRepository[T] {
	return &PostgresRepository[T]{db: db}
}

func (r *PostgresRepository[T]) DB() *pgxpool.Pool { return r.db }

func (r *PostgresRepository[T]) Save(ctx context.Context, query string, args ...any) error {
	_, err := r.querier(ctx).Exec(ctx, query, args...)
	return err
}

// SaveWithResult behaves like Save but also returns the driver's CommandTag,
// so callers that need to inspect RowsAffected() (e.g. optimistic-concurrency
// upserts) can do so without the generic Save/Delete helpers having to change
// shape for every other repository.
func (r *PostgresRepository[T]) SaveWithResult(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return r.querier(ctx).Exec(ctx, query, args...)
}

func (r *PostgresRepository[T]) Delete(ctx context.Context, query string, args ...any) error {
	_, err := r.querier(ctx).Exec(ctx, query, args...)
	return err
}

func (r *PostgresRepository[T]) Find(ctx context.Context, query string, fromRow func(pgx.Row) (T, error), args ...any) (T, error) {
	row := r.querier(ctx).QueryRow(ctx, query, args...)
	return fromRow(row)
}

func (r *PostgresRepository[T]) FindMany(ctx context.Context, query string, fromRow func(pgx.Row) (T, error), args ...any) ([]T, error) {
	rows, err := r.querier(ctx).Query(ctx, query, args...)
	if err != nil {
		var zero []T
		return zero, fmt.Errorf("error fetching rows: %w", err)
	}
	defer rows.Close()
	results := make([]T, 0)
	for rows.Next() {
		result, err := fromRow(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r *PostgresRepository[T]) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	if err = fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}
