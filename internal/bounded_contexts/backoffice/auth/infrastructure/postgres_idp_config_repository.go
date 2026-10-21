package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresIdpConfigRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresIdpConfigRepository(pool *pgxpool.Pool) *PostgresIdpConfigRepository {
	return &PostgresIdpConfigRepository{pool: pool}
}

func (r *PostgresIdpConfigRepository) FindByReferenceName(ctx context.Context, referenceName vos.IdpReferenceName) (*auth.IdpConfig, error) {
	var name string

	err := r.pool.QueryRow(ctx, `
		SELECT reference_name
		FROM backoffice.identity_providers
		WHERE reference_name = $1 AND enabled = TRUE
	`, referenceName).Scan(&name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying identity provider: %w", err)
	}

	return &auth.IdpConfig{IdpReferenceName: name}, nil
}
