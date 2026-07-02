package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/domain/criteria"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/admin_assertion_keys"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAdminAssertionKeyRepository struct {
	postgres *infrastructure.PostgresRepository[*admin_assertion_keys.AdminAssertionKey]
	pool     *pgxpool.Pool
}

func NewPostgresAdminAssertionKeyRepository(db *pgxpool.Pool) *PostgresAdminAssertionKeyRepository {
	return &PostgresAdminAssertionKeyRepository{
		postgres: infrastructure.NewPostgresRepository[*admin_assertion_keys.AdminAssertionKey](db),
		pool:     db,
	}
}

func (r *PostgresAdminAssertionKeyRepository) Save(ctx context.Context, aak *admin_assertion_keys.AdminAssertionKey) error {
	p := aak.ToPrimitives()
	query := `
		INSERT INTO admin_assertion_keys (id, admin_id, algorithm, public_key)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			admin_id = EXCLUDED.admin_id,
			algorithm = EXCLUDED.algorithm,
			public_key = EXCLUDED.public_key
	`
	return r.postgres.Save(ctx, query, p.Id, p.AdminId, p.Algorithm, p.PublicKey)
}

func (r *PostgresAdminAssertionKeyRepository) SearchById(ctx context.Context, id ids2.AdminAssertionKeyId) (*admin_assertion_keys.AdminAssertionKey, error) {
	query := `SELECT id, admin_id, algorithm, public_key FROM admin_assertion_keys WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresAdminAssertionKeyRepository) FindByAdminId(ctx context.Context, adminID ids2.AdminId) (*admin_assertion_keys.AdminAssertionKey, error) {
	query := `SELECT id, admin_id, algorithm, public_key FROM admin_assertion_keys WHERE admin_id = $1 AND deleted_at IS NULL LIMIT 1`
	return r.postgres.Find(ctx, query, r.fromRow, adminID.String())
}

func (r *PostgresAdminAssertionKeyRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*admin_assertion_keys.AdminAssertionKey, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresAdminAssertionKeyRepository) SoftDelete(ctx context.Context, id ids2.AdminAssertionKeyId) error {
	query := `UPDATE admin_assertion_keys SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresAdminAssertionKeyRepository) HardDelete(ctx context.Context, id ids2.AdminAssertionKeyId) error {
	query := `DELETE FROM admin_assertion_keys WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresAdminAssertionKeyRepository) fromRow(row pgx.Row) (*admin_assertion_keys.AdminAssertionKey, error) {
	var p admin_assertion_keys.AdminAssertionKeyPrimitives
	err := row.Scan(&p.Id, &p.AdminId, &p.Algorithm, &p.PublicKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning admin assertion key: %w", err)
	}
	return admin_assertion_keys.FromPrimitives(p)
}
