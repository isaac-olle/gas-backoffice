package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTenantRepository struct {
	postgres *infrastructure.PostgresRepository[*tenants.Tenant]
}

func NewPostgresTenantRepository(db *pgxpool.Pool) *PostgresTenantRepository {
	return &PostgresTenantRepository{
		postgres: infrastructure.NewPostgresRepository[*tenants.Tenant](db),
	}
}

func (r *PostgresTenantRepository) Save(ctx context.Context, tenant *tenants.Tenant) error {
	p := tenant.ToPrimitives()
	query := `
		INSERT INTO tenants (id, name, admin_id, client_ids, protected_resource_ids, is_system)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			client_ids = EXCLUDED.client_ids,
			protected_resource_ids = EXCLUDED.protected_resource_ids
	`
	return r.postgres.Save(ctx, query, p.Id, p.Name, p.AdminId, p.ClientIds, p.ProtectedResourceIds, p.IsSystem)
}

func (r *PostgresTenantRepository) SearchById(ctx context.Context, id ids.TenantId) (*tenants.Tenant, error) {
	query := `SELECT id, name, admin_id, client_ids, protected_resource_ids, is_system FROM tenants WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresTenantRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*tenants.Tenant, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresTenantRepository) SoftDelete(ctx context.Context, id ids.TenantId) error {
	query := `UPDATE tenants SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresTenantRepository) HardDelete(ctx context.Context, id ids.TenantId) error {
	query := `DELETE FROM tenants WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresTenantRepository) fromRow(row pgx.Row) (*tenants.Tenant, error) {
	var p tenants.TenantPrimitives
	err := row.Scan(&p.Id, &p.Name, &p.AdminId, &p.ClientIds, &p.ProtectedResourceIds, &p.IsSystem)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning tenant: %w", err)
	}
	return tenants.FromPrimitives(p)
}
