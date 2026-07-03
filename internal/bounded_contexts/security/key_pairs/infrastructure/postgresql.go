package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresKeyPairRepository struct {
	postgres *infrastructure.PostgresRepository[*key_pairs.KeyPair]
}

func NewPostgresKeyPairRepository(db *pgxpool.Pool) *PostgresKeyPairRepository {
	return &PostgresKeyPairRepository{
		postgres: infrastructure.NewPostgresRepository[*key_pairs.KeyPair](db),
	}
}

func (r *PostgresKeyPairRepository) Save(ctx context.Context, keyPair *key_pairs.KeyPair) error {
	p := keyPair.ToPrimitives()
	query := `
		INSERT INTO key_pairs (id, protected_resource_id, algorithm, public_key, private_key)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			protected_resource_id = EXCLUDED.protected_resource_id,
			algorithm = EXCLUDED.algorithm,
			public_key = EXCLUDED.public_key,
			private_key = EXCLUDED.private_key
	`
	return r.postgres.Save(ctx, query, p.Id, p.ProtectedResourceId, p.Algorithm, p.PublicKey, p.PrivateKey)
}

func (r *PostgresKeyPairRepository) SearchById(ctx context.Context, id ids.KeyPairId) (*key_pairs.KeyPair, error) {
	query := `SELECT id, protected_resource_id, algorithm, public_key, private_key FROM key_pairs WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresKeyPairRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*key_pairs.KeyPair, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresKeyPairRepository) SoftDelete(ctx context.Context, id ids.KeyPairId) error {
	query := `UPDATE key_pairs SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresKeyPairRepository) HardDelete(ctx context.Context, id ids.KeyPairId) error {
	query := `DELETE FROM key_pairs WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresKeyPairRepository) fromRow(row pgx.Row) (*key_pairs.KeyPair, error) {
	var p key_pairs.KeyPairPrimitives
	err := row.Scan(&p.Id, &p.ProtectedResourceId, &p.Algorithm, &p.PublicKey, &p.PrivateKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning key pair: %w", err)
	}
	return key_pairs.FromPrimitives(p)
}
