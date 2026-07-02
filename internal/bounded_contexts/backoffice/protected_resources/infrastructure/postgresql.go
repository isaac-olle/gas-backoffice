package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources"
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
	query := `
		INSERT INTO protected_resources (id, name, tenant_id, admin_id, scope_ids, is_system)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			scope_ids = EXCLUDED.scope_ids
	`
	return r.postgres.Save(ctx, query, p.Id, p.Name, p.TenantId, p.AdminId, p.ScopeIds, p.IsSystem)
}

func (r *PostgresProtectedResourceRepository) SearchById(ctx context.Context, id ids.ProtectedResourceId) (*protected_resources.ProtectedResource, error) {
	query := `SELECT id, name, tenant_id, admin_id, scope_ids, is_system FROM protected_resources WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresProtectedResourceRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*protected_resources.ProtectedResource, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresProtectedResourceRepository) SoftDelete(ctx context.Context, id ids.ProtectedResourceId) error {
	query := `UPDATE protected_resources SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresProtectedResourceRepository) HardDelete(ctx context.Context, id ids.ProtectedResourceId) error {
	query := `DELETE FROM protected_resources WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresProtectedResourceRepository) fromRow(row pgx.Row) (*protected_resources.ProtectedResource, error) {
	var p protected_resources.ProtectedResourcePrimitives
	err := row.Scan(&p.Id, &p.Name, &p.TenantId, &p.AdminId, &p.ScopeIds, &p.IsSystem)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning protected resource: %w", err)
	}
	return protected_resources.FromPrimitives(p)
}
