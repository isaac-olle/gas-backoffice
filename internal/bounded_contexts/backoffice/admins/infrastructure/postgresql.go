package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
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
		INSERT INTO backoffice.users (id, name, email, role, is_system, version, subscription_ids)
		VALUES ($1, $2, $3, 'admin', $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			subscription_ids = EXCLUDED.subscription_ids,
			version = EXCLUDED.version + 1
		WHERE backoffice.users.version = EXCLUDED.version
	`
	tag, err := r.postgres.SaveWithResult(ctx, query, p.Id, p.Name, p.Email, p.IsSystem, p.Version, p.SubscriptionIds)
	if err != nil {
		return fmt.Errorf("error saving admin: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return persistance.ErrConcurrentModification
	}
	return nil
}

func (r *PostgresAdminRepository) SearchById(ctx context.Context, id ids.UserId) (*admins.Admin, error) {
	query := `SELECT id, name, email, is_system, version, subscription_ids FROM backoffice.users WHERE role = 'admin' AND id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresAdminRepository) Matching(ctx context.Context, c criteria.Criteria) ([]*admins.Admin, error) {
	where, args, err := c.ToPostgresql(1)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, name, email, is_system, version, subscription_ids FROM backoffice.users WHERE role = 'admin' AND deleted_at IS NULL`
	if where != "" {
		query += " AND " + where
	}
	return r.postgres.FindMany(ctx, query, r.fromRow, args...)
}

func (r *PostgresAdminRepository) SoftDelete(ctx context.Context, id ids.UserId) error {
	query := `UPDATE backoffice.users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresAdminRepository) HardDelete(ctx context.Context, id ids.UserId) error {
	query := `DELETE FROM backoffice.users WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresAdminRepository) fromRow(row pgx.Row) (*admins.Admin, error) {
	var p admins.AdminPrimitives
	var email *string
	err := row.Scan(&p.Id, &p.Name, &email, &p.IsSystem, &p.Version, &p.SubscriptionIds)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning admin: %w", err)
	}
	if email != nil {
		p.Email = *email
	}
	return admins.FromPrimitives(p)
}
