package infrastructure

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresSigningKeyRepository reads the public half of the key pairs that
// sign tokens. It joins security.key_pairs with backoffice.protected_resources
// so keys of disabled or deleted protected resources are not published. It
// never touches the private_key column.
type PostgresSigningKeyRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresSigningKeyRepository(pool *pgxpool.Pool) *PostgresSigningKeyRepository {
	return &PostgresSigningKeyRepository{pool: pool}
}

func (r *PostgresSigningKeyRepository) SearchActive(ctx context.Context) ([]*jwks.SigningKeyProjection, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT kp.id, kp.algorithm, kp.public_key
		FROM security.key_pairs kp
		JOIN backoffice.protected_resources pr ON pr.id = kp.protected_resource_id
		WHERE kp.deleted_at IS NULL
		  AND pr.deleted_at IS NULL
		  AND pr.enabled = TRUE
		ORDER BY kp.created_at
	`)
	if err != nil {
		return nil, fmt.Errorf("error querying signing keys: %w", err)
	}
	defer rows.Close()

	keys := make([]*jwks.SigningKeyProjection, 0)
	for rows.Next() {
		var kid, algorithm string
		var publicKey []byte
		if err = rows.Scan(&kid, &algorithm, &publicKey); err != nil {
			return nil, fmt.Errorf("error scanning signing key: %w", err)
		}
		keys = append(keys, jwks.NewSigningKeyProjection(kid, algorithm, publicKey))
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating signing keys: %w", err)
	}
	return keys, nil
}
