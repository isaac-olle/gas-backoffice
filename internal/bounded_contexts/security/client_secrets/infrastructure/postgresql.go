package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
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
	authorId, authorRole := authorPrimitives(ctx)
	query := `
		INSERT INTO security.client_secrets (id, client_id, secret_hash, author_id, author_role)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			client_id = EXCLUDED.client_id,
			secret_hash = EXCLUDED.secret_hash
	`
	return r.postgres.Save(ctx, query, p.Id, p.ClientId, p.SecretHash, authorId, authorRole)
}

func (r *PostgresClientSecretRepository) SearchById(ctx context.Context, id ids2.ClientSecretId) (*client_secrets.ClientSecret, error) {
	query := `SELECT id, client_id, secret_hash FROM security.client_secrets WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresClientSecretRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*client_secrets.ClientSecret, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresClientSecretRepository) SoftDelete(ctx context.Context, id ids2.ClientSecretId) error {
	query := `UPDATE security.client_secrets SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresClientSecretRepository) HardDelete(ctx context.Context, id ids2.ClientSecretId) error {
	query := `DELETE FROM security.client_secrets WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresClientSecretRepository) FindByClientId(ctx context.Context, clientID ids2.ClientId) (*client_secrets.ClientSecret, error) {
	query := `SELECT id, client_id, secret_hash FROM security.client_secrets WHERE client_id = $1 AND deleted_at IS NULL LIMIT 1`
	return r.postgres.Find(ctx, query, r.fromRow, clientID.String())
}

func authorPrimitives(ctx context.Context) (string, string) {
	author, ok := vo.AuthorFromContext[vo.Author](ctx)
	if !ok {
		return "", ""
	}
	p := author.ToPrimitives()
	return p.Id, p.Role
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
