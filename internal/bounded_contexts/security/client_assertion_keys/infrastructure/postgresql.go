package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClientAssertionKeyRepository struct {
	postgres *infrastructure.PostgresRepository[*client_assertion_keys.ClientAssertionKey]
}

func NewPostgresClientAssertionKeyRepository(db *pgxpool.Pool) *PostgresClientAssertionKeyRepository {
	return &PostgresClientAssertionKeyRepository{
		postgres: infrastructure.NewPostgresRepository[*client_assertion_keys.ClientAssertionKey](db),
	}
}

func (r *PostgresClientAssertionKeyRepository) Save(ctx context.Context, cak *client_assertion_keys.ClientAssertionKey) error {
	p := cak.ToPrimitives()
	authorId, authorRole := actorPrimitives(ctx)
	query := `
		INSERT INTO client_assertion_keys (id, client_id, algorithm, public_key, author_id, author_role)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			client_id = EXCLUDED.client_id,
			algorithm = EXCLUDED.algorithm,
			public_key = EXCLUDED.public_key
	`
	return r.postgres.Save(ctx, query, p.Id, p.ClientId, p.Algorithm, p.PublicKey, authorId, authorRole)
}

func (r *PostgresClientAssertionKeyRepository) SearchById(ctx context.Context, id ids2.ClientAssertionKeyId) (*client_assertion_keys.ClientAssertionKey, error) {
	query := `SELECT id, client_id, key_pair_id, algorithm, public_key FROM client_assertion_keys WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresClientAssertionKeyRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*client_assertion_keys.ClientAssertionKey, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresClientAssertionKeyRepository) SoftDelete(ctx context.Context, id ids2.ClientAssertionKeyId) error {
	query := `UPDATE client_assertion_keys SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresClientAssertionKeyRepository) HardDelete(ctx context.Context, id ids2.ClientAssertionKeyId) error {
	query := `DELETE FROM client_assertion_keys WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresClientAssertionKeyRepository) FindByClientId(ctx context.Context, clientID ids2.ClientId) (*client_assertion_keys.ClientAssertionKey, error) {
	query := `SELECT id, client_id, key_pair_id, algorithm, public_key FROM client_assertion_keys WHERE client_id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, clientID.String())
}

func actorPrimitives(ctx context.Context) (string, string) {
	actor, ok := vo.ActorFromContext[vo.Actor](ctx)
	if !ok {
		return "", ""
	}
	p := actor.ToPrimitives()
	return p.Id, p.Role
}

func (r *PostgresClientAssertionKeyRepository) fromRow(row pgx.Row) (*client_assertion_keys.ClientAssertionKey, error) {
	var p client_assertion_keys.ClientAssertionKeyPrimitives
	err := row.Scan(&p.Id, &p.ClientId, &p.Algorithm, &p.PublicKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning client assertion key: %w", err)
	}
	return client_assertion_keys.FromPrimitives(p)
}
