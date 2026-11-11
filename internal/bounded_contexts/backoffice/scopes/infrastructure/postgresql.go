package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresScopeRepository struct {
	postgres *infrastructure.PostgresRepository[*scopes.Scope]
}

func NewPostgresScopeRepository(db *pgxpool.Pool) *PostgresScopeRepository {
	return &PostgresScopeRepository{
		postgres: infrastructure.NewPostgresRepository[*scopes.Scope](db),
	}
}

func (r *PostgresScopeRepository) Save(ctx context.Context, scope *scopes.Scope) error {
	p := scope.ToPrimitives()
	authorId, authorRole := authorPrimitives(ctx)
	query := `
		INSERT INTO backoffice.scopes (id, name, protected_resource_id, subscription_id, allowed_client_ids, is_system, enabled, author_id, author_role, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			allowed_client_ids = EXCLUDED.allowed_client_ids,
			enabled = EXCLUDED.enabled,
			version = EXCLUDED.version + 1
		WHERE backoffice.scopes.version = EXCLUDED.version
	`
	tag, err := r.postgres.SaveWithResult(ctx, query,
		p.Id, p.Name, p.ProtectedResourceId, p.SubscriptionId, p.AllowedClientIds, p.IsSystem, p.Enabled, authorId, authorRole, p.Version)
	if err != nil {
		return fmt.Errorf("error saving scope: %w", err)
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

func (r *PostgresScopeRepository) SearchById(ctx context.Context, id ids.ScopeId) (*scopes.Scope, error) {
	query := `SELECT id, name, protected_resource_id, subscription_id, allowed_client_ids, is_system, enabled, version, author_id, author_role
		FROM backoffice.scopes WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresScopeRepository) Matching(ctx context.Context, c criteria.Criteria) ([]*scopes.Scope, error) {
	where, args, err := c.ToPostgresql(1)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, name, protected_resource_id, subscription_id, allowed_client_ids, is_system, enabled, version, author_id, author_role FROM backoffice.scopes WHERE deleted_at IS NULL`
	if where != "" {
		query += " AND " + where
	}
	return r.postgres.FindMany(ctx, query, r.fromRow, args...)
}

func (r *PostgresScopeRepository) FindByAllowedClientId(ctx context.Context, clientId ids.ClientId) ([]*scopes.Scope, error) {
	query := `SELECT id, name, protected_resource_id, subscription_id, allowed_client_ids, is_system, enabled, version, author_id, author_role
		FROM backoffice.scopes WHERE $1 = ANY(allowed_client_ids) AND deleted_at IS NULL`
	return r.postgres.FindMany(ctx, query, r.fromRow, clientId.String())
}

func (r *PostgresScopeRepository) SoftDelete(ctx context.Context, id ids.ScopeId) error {
	query := `UPDATE backoffice.scopes SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresScopeRepository) HardDelete(ctx context.Context, id ids.ScopeId) error {
	query := `DELETE FROM backoffice.scopes WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresScopeRepository) fromRow(row pgx.Row) (*scopes.Scope, error) {
	var p scopes.ScopePrimitives
	err := row.Scan(&p.Id, &p.Name, &p.ProtectedResourceId, &p.SubscriptionId, &p.AllowedClientIds, &p.IsSystem, &p.Enabled, &p.Version, &p.AuthorId, &p.AuthorRole)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, app_errors.NotFoundError(scopes.ErrScopeNotFound).WithDescription("scope not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning scope: %w", err)
	}
	return scopes.FromPrimitives(p)
}
