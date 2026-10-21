package seed

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAlreadySeeded = errors.New("already seeded")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) IsSeeded(ctx context.Context) (bool, error) {
	var isSeeded bool
	err := r.pool.QueryRow(ctx, `SELECT is_seed FROM public.seed`).Scan(&isSeeded)
	return isSeeded, err
}

func (r *Repository) MarkSeeded(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `UPDATE public.seed SET is_seed = TRUE`)
	if err != nil {
		return fmt.Errorf("mark seeded: %w", err)
	}
	return nil
}
