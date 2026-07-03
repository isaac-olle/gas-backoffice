package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClientSecretRepository struct {
	postgres *infrastructure.PostgresRepository[*client_secrets.ClientSecret]
}

func NewPostgresClientSecretRepository(db *pgxpool.Pool) *PostgresClientSecretRepository {
	return &PostgresClientSecretRepository{
		postgres: infrastructure.NewPostgresRepository[*client_secrets.ClientSecret](db),
	}
}

func (r *PostgresClientSecretRepository) Save(ctx context.Context, cs *client_secrets.ClientSecret) error {
	p := cs.ToPrimitives()
	query := `
		INSERT INTO client_secrets (id, client_id, secret_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET
			client_id = EXCLUDED.client_id,
			secret_hash = EXCLUDED.secret_hash
	`
	return r.postgres.Save(ctx, query, p.Id, p.ClientId, p.SecretHash)
}

func (r *PostgresClientSecretRepository) SearchById(ctx context.Context, id ids2.ClientSecretId) (*client_secrets.ClientSecret, error) {
	query := `SELECT id, client_id, secret_hash FROM client_secrets WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresClientSecretRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*client_secrets.ClientSecret, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresClientSecretRepository) SoftDelete(ctx context.Context, id ids2.ClientSecretId) error {
	query := `UPDATE client_secrets SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresClientSecretRepository) HardDelete(ctx context.Context, id ids2.ClientSecretId) error {
	query := `DELETE FROM client_secrets WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresClientSecretRepository) FindByClientId(ctx context.Context, clientID ids2.ClientId) (*client_secrets.ClientSecret, error) {
	query := `SELECT id, client_id, secret_hash FROM client_secrets WHERE client_id = $1 AND deleted_at IS NULL LIMIT 1`
	return r.postgres.Find(ctx, query, r.fromRow, clientID.String())
}

func (r *PostgresClientSecretRepository) fromRow(row pgx.Row) (*client_secrets.ClientSecret, error) {
	var p client_secrets.ClientSecretPrimitives
	err := row.Scan(&p.Id, &p.ClientId, &p.SecretHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, app_errors.NotFoundError(client_secrets.ErrClientSecretNotFound).WithDescription("client secret not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning client secret: %w", err)
	}
	return client_secrets.FromPrimitives(p)
}
