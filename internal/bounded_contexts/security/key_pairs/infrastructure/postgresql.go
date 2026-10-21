package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresKeyPairRepository struct {
	postgres  *infrastructure.PostgresRepository[*key_pairs.KeyPair]
	encrypter persistance.Encrypter
}

func NewPostgresKeyPairRepository(db *pgxpool.Pool, encrypter persistance.Encrypter) *PostgresKeyPairRepository {
	return &PostgresKeyPairRepository{
		postgres:  infrastructure.NewPostgresRepository[*key_pairs.KeyPair](db),
		encrypter: encrypter,
	}
}

func (r *PostgresKeyPairRepository) Save(ctx context.Context, keyPair *key_pairs.KeyPair) error {
	p := keyPair.ToPrimitives()
	authorId, authorRole := authorPrimitives(ctx)

	encryptedPrivateKey, err := r.encrypter.Encrypt(string(p.PrivateKey))
	if err != nil {
		return fmt.Errorf("error encrypting private key: %w", err)
	}

	query := `
		INSERT INTO security.key_pairs (id, protected_resource_id, algorithm, public_key, private_key, author_id, author_role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			protected_resource_id = EXCLUDED.protected_resource_id,
			algorithm = EXCLUDED.algorithm,
			public_key = EXCLUDED.public_key,
			private_key = EXCLUDED.private_key
	`
	return r.postgres.Save(ctx, query, p.Id, p.ProtectedResourceId, p.Algorithm, p.PublicKey, []byte(encryptedPrivateKey), authorId, authorRole)
}

func (r *PostgresKeyPairRepository) SearchById(ctx context.Context, id ids.KeyPairId) (*key_pairs.KeyPair, error) {
	query := `SELECT id, protected_resource_id, algorithm, public_key, private_key FROM security.key_pairs WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresKeyPairRepository) SearchByProtectedResourceId(ctx context.Context, id ids.ProtectedResourceId) (*key_pairs.KeyPair, error) {
	query := `SELECT id, protected_resource_id, algorithm, public_key, private_key FROM security.key_pairs WHERE protected_resource_id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresKeyPairRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*key_pairs.KeyPair, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresKeyPairRepository) SoftDelete(ctx context.Context, id ids.KeyPairId) error {
	query := `UPDATE security.key_pairs SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresKeyPairRepository) HardDelete(ctx context.Context, id ids.KeyPairId) error {
	query := `DELETE FROM security.key_pairs WHERE id = $1`
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

func (r *PostgresKeyPairRepository) fromRow(row pgx.Row) (*key_pairs.KeyPair, error) {
	var p key_pairs.KeyPairPrimitives
	var encryptedPrivateKey []byte
	err := row.Scan(&p.Id, &p.ProtectedResourceId, &p.Algorithm, &p.PublicKey, &encryptedPrivateKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning key pair: %w", err)
	}

	privateKey, err := r.encrypter.Decrypt(string(encryptedPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("error decrypting private key: %w", err)
	}
	p.PrivateKey = []byte(privateKey)

	return key_pairs.FromPrimitives(p)
}
