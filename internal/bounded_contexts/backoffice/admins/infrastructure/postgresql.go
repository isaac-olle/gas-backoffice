package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAdminRepository struct {
	postgres *infrastructure.PostgresRepository[*admins.Admin]
}

func NewPostgresAdminRepository(db *pgxpool.Pool) *PostgresAdminRepository {
	return &PostgresAdminRepository{
		postgres: infrastructure.NewPostgresRepository[*admins.Admin](db),
	}
}

func (r *PostgresAdminRepository) Save(ctx context.Context, admin *admins.Admin) error {
	p := admin.ToPrimitives()
	query := `
		INSERT INTO admins (id, name, tenant_ids, is_system)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			tenant_ids = EXCLUDED.tenant_ids
	`
	return r.postgres.Save(ctx, query, p.Id, p.Name, p.TenantIds, p.IsSystem)
}

func (r *PostgresAdminRepository) SearchById(ctx context.Context, id ids.AdminId) (*admins.Admin, error) {
	query := `SELECT id, name, tenant_ids, is_system FROM admins WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresAdminRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*admins.Admin, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresAdminRepository) SoftDelete(ctx context.Context, id ids.AdminId) error {
	query := `UPDATE admins SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresAdminRepository) HardDelete(ctx context.Context, id ids.AdminId) error {
	query := `DELETE FROM admins WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresAdminRepository) fromRow(row pgx.Row) (*admins.Admin, error) {
	var p admins.AdminPrimitives
	err := row.Scan(&p.Id, &p.Name, &p.TenantIds, &p.IsSystem)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, app_errors.NotFoundError(admins.ErrAdminNotFound).WithDescription("admin not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning admin: %w", err)
	}
	return admins.FromPrimitives(p)
}
