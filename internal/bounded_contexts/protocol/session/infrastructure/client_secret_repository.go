package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/app_errors"

	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSessionClientSecretRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresSessionClientSecretRepository(pool *pgxpool.Pool) *PostgresSessionClientSecretRepository {
	return &PostgresSessionClientSecretRepository{pool: pool}
}

func (r *PostgresSessionClientSecretRepository) FindHashByClientId(ctx context.Context, clientID string) (string, error) {
	var hash string
	err := r.pool.QueryRow(ctx,
		`SELECT secret_hash FROM client_secrets WHERE client_id = $1 AND deleted_at IS NULL LIMIT 1`,
		clientID,
	).Scan(&hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", app_errors.NotFoundError(session.ErrClientSecretNotFound).
			WithDescription("no client secret found for client")
	}
	if err != nil {
		return "", fmt.Errorf("query client secret hash: %w", err)
	}
	return hash, nil
}
