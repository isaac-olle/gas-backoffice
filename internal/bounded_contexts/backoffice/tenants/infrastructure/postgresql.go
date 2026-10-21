package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
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

// Save upserts both the shared backoffice.users row and its
// backoffice.tenant_profiles extension row in the same transaction, keeping
// the existing OCC version check that already lives on backoffice.users.
func (r *PostgresTenantRepository) Save(ctx context.Context, tenant *tenants.Tenant) error {
	p := tenant.ToPrimitives()
	userQuery := `
		INSERT INTO backoffice.users (id, name, email, role, is_system, version, subscription_ids)
		VALUES ($1, $2, $3, 'tenant', $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			subscription_ids = EXCLUDED.subscription_ids,
			version = EXCLUDED.version + 1
		WHERE backoffice.users.version = EXCLUDED.version
	`
	profileQuery := `
		INSERT INTO backoffice.tenant_profiles (user_id, admin_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET admin_id = EXCLUDED.admin_id
	`

	return r.postgres.WithTx(ctx, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, userQuery, p.Id, p.Name, p.Email, p.IsSystem, p.Version, p.SubscriptionIds)
		if err != nil {
			return fmt.Errorf("error saving tenant: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return persistance.ErrConcurrentModification
		}
		if _, err = tx.Exec(ctx, profileQuery, p.Id, p.AdminId); err != nil {
			return fmt.Errorf("error saving tenant profile: %w", err)
		}
		return nil
	})
}

func (r *PostgresTenantRepository) SearchById(ctx context.Context, id ids.UserId) (*tenants.Tenant, error) {
	query := `
		SELECT u.id, u.name, u.email, u.is_system, u.version, tp.admin_id, u.subscription_ids
		FROM backoffice.users u
		JOIN backoffice.tenant_profiles tp ON tp.user_id = u.id
		WHERE u.role = 'tenant' AND u.id = $1 AND u.deleted_at IS NULL
	`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresTenantRepository) Matching(ctx context.Context, c criteria.Criteria) ([]*tenants.Tenant, error) {
	where, args, err := c.ToPostgresql(1)
	if err != nil {
		return nil, err
	}
	query := `
		SELECT u.id, u.name, u.email, u.is_system, u.version, tp.admin_id, u.subscription_ids
		FROM backoffice.users u
		JOIN backoffice.tenant_profiles tp ON tp.user_id = u.id
		WHERE u.role = 'tenant' AND u.deleted_at IS NULL
	`
	if where != "" {
		query += " AND " + where
	}
	return r.postgres.FindMany(ctx, query, r.fromRow, args...)
}

func (r *PostgresTenantRepository) SoftDelete(ctx context.Context, id ids.UserId) error {
	query := `UPDATE backoffice.users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

// HardDelete removes the tenant_profiles extension row before the shared
// users row, since the former has a foreign key on the latter.
func (r *PostgresTenantRepository) HardDelete(ctx context.Context, id ids.UserId) error {
	return r.postgres.WithTx(ctx, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `DELETE FROM backoffice.tenant_profiles WHERE user_id = $1`, id.String()); err != nil {
			return fmt.Errorf("error deleting tenant profile: %w", err)
		}
		if _, err := tx.Exec(ctx, `DELETE FROM backoffice.users WHERE id = $1`, id.String()); err != nil {
			return fmt.Errorf("error deleting tenant: %w", err)
		}
		return nil
	})
}

func (r *PostgresTenantRepository) fromRow(row pgx.Row) (*tenants.Tenant, error) {
	var p tenants.TenantPrimitives
	var email *string
	err := row.Scan(&p.Id, &p.Name, &email, &p.IsSystem, &p.Version, &p.AdminId, &p.SubscriptionIds)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning tenant: %w", err)
	}
	if email != nil {
		p.Email = *email
	}
	return tenants.FromPrimitives(p)
}
