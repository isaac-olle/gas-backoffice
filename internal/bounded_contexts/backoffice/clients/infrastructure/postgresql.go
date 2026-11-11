package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClientRepository struct {
	postgres *infrastructure.PostgresRepository[*clients.Client]
}

func NewPostgresClientRepository(db *pgxpool.Pool) *PostgresClientRepository {
	return &PostgresClientRepository{
		postgres: infrastructure.NewPostgresRepository[*clients.Client](db),
	}
}

func (r *PostgresClientRepository) Save(ctx context.Context, client *clients.Client) error {
	p := client.ToPrimitives()
	authorId, authorRole := authorPrimitives(ctx)
	query := `
		INSERT INTO backoffice.clients (id, name, subscription_id, enabled, public, allowed_grant_types, redirect_uris, assertion_key_ids, secret_ids, allowed_scopes, author_id, author_role, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			enabled = EXCLUDED.enabled,
			public = EXCLUDED.public,
			allowed_grant_types = EXCLUDED.allowed_grant_types,
			redirect_uris = EXCLUDED.redirect_uris,
			assertion_key_ids = EXCLUDED.assertion_key_ids,
			secret_ids = EXCLUDED.secret_ids,
			allowed_scopes = EXCLUDED.allowed_scopes,
			version = EXCLUDED.version + 1
		WHERE backoffice.clients.version = EXCLUDED.version
	`
	tag, err := r.postgres.SaveWithResult(ctx, query,
		p.Id, p.Name, p.SubscriptionId, p.Enabled, p.Public, p.AllowedGrantTypes,
		p.RedirectUris, p.AssertionKeyIds, p.SecretIds, p.AllowedScopes, authorId, authorRole, p.Version)
	if err != nil {
		return fmt.Errorf("error saving client: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return persistance.ErrConcurrentModification
	}
	return nil
}

func authorPrimitives(ctx context.Context) (string, string) {
	author, ok := vo.AuthorFromContext(ctx)
	if !ok {
		return "", ""
	}
	p := author.ToPrimitives()
	return p.Id, p.Role
}

func (r *PostgresClientRepository) SearchById(ctx context.Context, id ids.ClientId) (*clients.Client, error) {
	query := `SELECT id, name, subscription_id, enabled, public, is_system,
		allowed_grant_types, redirect_uris, assertion_key_ids, secret_ids, allowed_scopes, version, author_id, author_role
		FROM backoffice.clients WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresClientRepository) Matching(ctx context.Context, c criteria.Criteria) ([]*clients.Client, error) {
	where, args, err := c.ToPostgresql(1)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, name, subscription_id, enabled, public, is_system,
		allowed_grant_types, redirect_uris, assertion_key_ids, secret_ids, allowed_scopes, version, author_id, author_role
		FROM backoffice.clients WHERE deleted_at IS NULL`
	if where != "" {
		query += " AND " + where
	}
	return r.postgres.FindMany(ctx, query, r.fromRow, args...)
}

func (r *PostgresClientRepository) SoftDelete(ctx context.Context, id ids.ClientId) error {
	query := `UPDATE backoffice.clients SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresClientRepository) HardDelete(ctx context.Context, id ids.ClientId) error {
	query := `DELETE FROM backoffice.clients WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresClientRepository) fromRow(row pgx.Row) (*clients.Client, error) {
	var p clients.ClientPrimitives
	err := row.Scan(&p.Id,
		&p.Name,
		&p.SubscriptionId,
		&p.Enabled,
		&p.Public,
		&p.IsSystem,
		&p.AllowedGrantTypes,
		&p.RedirectUris,
		&p.AssertionKeyIds,
		&p.SecretIds,
		&p.AllowedScopes,
		&p.Version,
		&p.AuthorId,
		&p.AuthorRole)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning client: %w", err)
	}
	return clients.FromPrimitives(p)
}
