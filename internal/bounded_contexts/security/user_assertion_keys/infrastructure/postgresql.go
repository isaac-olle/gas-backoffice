package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/user_assertion_keys"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserAssertionKeyRepository struct {
	postgres *infrastructure.PostgresRepository[*user_assertion_keys.UserAssertionKey]
	pool     *pgxpool.Pool
}

func NewPostgresUserAssertionKeyRepository(db *pgxpool.Pool) *PostgresUserAssertionKeyRepository {
	return &PostgresUserAssertionKeyRepository{
		postgres: infrastructure.NewPostgresRepository[*user_assertion_keys.UserAssertionKey](db),
		pool:     db,
	}
}

func (r *PostgresUserAssertionKeyRepository) Save(ctx context.Context, uak *user_assertion_keys.UserAssertionKey) error {
	p := uak.ToPrimitives()
	authorId, authorRole := authorPrimitives(ctx)
	query := `
		INSERT INTO security.user_assertion_keys (id, user_id, algorithm, public_key, author_id, author_role)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			algorithm = EXCLUDED.algorithm,
			public_key = EXCLUDED.public_key
	`
	return r.postgres.Save(ctx, query, p.Id, p.UserId, p.Algorithm, p.PublicKey, authorId, authorRole)
}

func (r *PostgresUserAssertionKeyRepository) SearchById(ctx context.Context, id ids2.UserAssertionKeyId) (*user_assertion_keys.UserAssertionKey, error) {
	query := `SELECT id, user_id, algorithm, public_key FROM security.user_assertion_keys WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresUserAssertionKeyRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*user_assertion_keys.UserAssertionKey, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresUserAssertionKeyRepository) SoftDelete(ctx context.Context, id ids2.UserAssertionKeyId) error {
	query := `UPDATE security.user_assertion_keys SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresUserAssertionKeyRepository) HardDelete(ctx context.Context, id ids2.UserAssertionKeyId) error {
	query := `DELETE FROM security.user_assertion_keys WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func authorPrimitives(ctx context.Context) (string, string) {
	author, ok := vo.AuthorFromContext[vo.Author](ctx)
	if !ok {
		return "", ""
	}
	p := author.ToPrimitives()
	return p.Id, p.Role
}

func (r *PostgresUserAssertionKeyRepository) fromRow(row pgx.Row) (*user_assertion_keys.UserAssertionKey, error) {
	var p user_assertion_keys.UserAssertionKeyPrimitives
	err := row.Scan(&p.Id, &p.UserId, &p.Algorithm, &p.PublicKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning user assertion key: %w", err)
	}
	return user_assertion_keys.FromPrimitives(p)
}
