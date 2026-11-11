package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	sharedInfrastructure "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSessionClientSecretRepository struct {
	postgres *sharedInfrastructure.PostgresRepository[*session.Client]
}

func NewPostgresSessionClientSecretRepository(pool *pgxpool.Pool) *PostgresSessionClientSecretRepository {
	return &PostgresSessionClientSecretRepository{postgres: sharedInfrastructure.NewPostgresRepository[*session.Client](pool)}
}

func (r *PostgresSessionClientSecretRepository) FindHashByClientId(ctx context.Context, clientID string) (string, error) {
	var hash string
	err := r.postgres.QueryRow(ctx,
		`SELECT secret_hash FROM security.client_secrets WHERE client_id = $1 AND deleted_at IS NULL LIMIT 1`,
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
