package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
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
	query := `
		INSERT INTO scopes (id, name, protected_resource_id, tenant_id, admin_id, allowed_client_ids, is_system)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			protected_resource_id = EXCLUDED.protected_resource_id,
			tenant_id = EXCLUDED.tenant_id,
			admin_id = EXCLUDED.admin_id,
			allowed_client_ids = EXCLUDED.allowed_client_ids
	`
	return r.postgres.Save(ctx, query, p.Id, p.Name, p.ProtectedResourceId, p.TenantId, p.AdminId, p.AllowedClientIds, p.IsSystem)
}

func (r *PostgresScopeRepository) SearchById(ctx context.Context, id ids.ScopeId) (*scopes.Scope, error) {
	query := `SELECT id, name, protected_resource_id, tenant_id, admin_id, allowed_client_ids, is_system FROM scopes WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresScopeRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*scopes.Scope, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresScopeRepository) SoftDelete(ctx context.Context, id ids.ScopeId) error {
	query := `UPDATE scopes SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresScopeRepository) HardDelete(ctx context.Context, id ids.ScopeId) error {
	query := `DELETE FROM scopes WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresScopeRepository) fromRow(row pgx.Row) (*scopes.Scope, error) {
	var p scopes.ScopePrimitives
	err := row.Scan(&p.Id, &p.Name, &p.ProtectedResourceId, &p.TenantId, &p.AdminId, &p.AllowedClientIds, &p.IsSystem)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, app_errors.NotFoundError(scopes.ErrScopeNotFound).WithDescription("scope not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning scope: %w", err)
	}
	return scopes.FromPrimitives(p)
}
