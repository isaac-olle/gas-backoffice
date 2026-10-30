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

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProtectedResourceRepository struct {
	postgres *infrastructure.PostgresRepository[*protected_resources.ProtectedResource]
}

func NewPostgresProtectedResourceRepository(db *pgxpool.Pool) *PostgresProtectedResourceRepository {
	return &PostgresProtectedResourceRepository{
		postgres: infrastructure.NewPostgresRepository[*protected_resources.ProtectedResource](db),
	}
}

func (r *PostgresProtectedResourceRepository) Save(ctx context.Context, pr *protected_resources.ProtectedResource) error {
	p := pr.ToPrimitives()
	authorId, authorRole := authorPrimitives(ctx)
	query := `
		INSERT INTO backoffice.protected_resources (id, name, subscription_id, is_system, enabled, scope_ids, author_id, author_role, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			enabled = EXCLUDED.enabled,
			scope_ids = EXCLUDED.scope_ids,
			version = EXCLUDED.version + 1
		WHERE backoffice.protected_resources.version = EXCLUDED.version
	`
	tag, err := r.postgres.SaveWithResult(ctx, query,
		p.Id, p.Name, p.SubscriptionId, p.IsSystem, p.Enabled, p.ScopeIds, authorId, authorRole, p.Version)
	if err != nil {
		return fmt.Errorf("error saving protected resource: %w", err)
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

func (r *PostgresProtectedResourceRepository) SearchById(ctx context.Context, id ids.ProtectedResourceId) (*protected_resources.ProtectedResource, error) {
	query := `SELECT id, name, subscription_id, is_system, enabled, scope_ids, version, author_id, author_role FROM backoffice.protected_resources WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresProtectedResourceRepository) Matching(ctx context.Context, c criteria.Criteria) ([]*protected_resources.ProtectedResource, error) {
	where, args, err := c.ToPostgresql(1)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, name, subscription_id, is_system, enabled, scope_ids, version, author_id, author_role FROM backoffice.protected_resources WHERE deleted_at IS NULL`
	if where != "" {
		query += " AND " + where
	}
	return r.postgres.FindMany(ctx, query, r.fromRow, args...)
}

func (r *PostgresProtectedResourceRepository) SoftDelete(ctx context.Context, id ids.ProtectedResourceId) error {
	query := `UPDATE backoffice.protected_resources SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresProtectedResourceRepository) HardDelete(ctx context.Context, id ids.ProtectedResourceId) error {
	query := `DELETE FROM backoffice.protected_resources WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresProtectedResourceRepository) fromRow(row pgx.Row) (*protected_resources.ProtectedResource, error) {
	var p protected_resources.ProtectedResourcePrimitives
	err := row.Scan(&p.Id, &p.Name, &p.SubscriptionId, &p.IsSystem, &p.Enabled, &p.ScopeIds, &p.Version, &p.AuthorId, &p.AuthorRole)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning protected resource: %w", err)
	}
	return protected_resources.FromPrimitives(p)
}
