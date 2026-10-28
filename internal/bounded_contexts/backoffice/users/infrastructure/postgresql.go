package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	postgres *infrastructure.PostgresRepository[*users.User]
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		postgres: infrastructure.NewPostgresRepository[*users.User](db),
	}
}

// Save upserts the shared backoffice.users row. User provenance is stored in
// users.author_id and users.author_role.
func (r *PostgresUserRepository) Save(ctx context.Context, user *users.User) error {
	p := user.ToPrimitives()
	authorId, authorRole := authorPrimitives(ctx)
	userQuery := `
		INSERT INTO backoffice.users (id, name, email, system_role, is_system, version, subscription_ids, memberships, author_id, author_role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			subscription_ids = EXCLUDED.subscription_ids,
			memberships = EXCLUDED.memberships,
			version = EXCLUDED.version + 1
		WHERE backoffice.users.version = EXCLUDED.version
	`
	memberships, err := json.Marshal(p.Memberships)
	if err != nil {
		return fmt.Errorf("marshal user memberships: %w", err)
	}
	return r.postgres.WithTx(ctx, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, userQuery, p.Id, p.Name, p.Email, p.SystemRole, p.IsSystem, p.Version, p.SubscriptionIds, memberships, authorId, authorRole)
		if err != nil {
			return fmt.Errorf("error saving user: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return persistance.ErrConcurrentModification
		}
		return nil
	})
}

func authorPrimitives(ctx context.Context) (string, string) {
	author, ok := vo.AuthorFromContext[vo.Author](ctx)
	if !ok {
		return "", ""
	}
	p := author.ToPrimitives()
	return p.Id, p.Role
}

func (r *PostgresUserRepository) SearchById(ctx context.Context, id ids.UserId) (*users.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.system_role, u.is_system, u.version, u.author_id, u.author_role, u.subscription_ids, u.memberships
		FROM backoffice.users u
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresUserRepository) Matching(ctx context.Context, c criteria.Criteria) ([]*users.User, error) {
	where, args, err := c.ToPostgresql(1)
	if err != nil {
		return nil, err
	}
	query := `
		SELECT u.id, u.name, u.email, u.system_role, u.is_system, u.version, u.author_id, u.author_role, u.subscription_ids, u.memberships
		FROM backoffice.users u
		WHERE u.deleted_at IS NULL
	`
	if where != "" {
		query += " AND " + where
	}
	return r.postgres.FindMany(ctx, query, r.fromRow, args...)
}

func (r *PostgresUserRepository) SoftDelete(ctx context.Context, id ids.UserId) error {
	query := `UPDATE backoffice.users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresUserRepository) HardDelete(ctx context.Context, id ids.UserId) error {
	return r.postgres.Delete(ctx, `DELETE FROM backoffice.users WHERE id = $1`, id.String())
}

func (r *PostgresUserRepository) fromRow(row pgx.Row) (*users.User, error) {
	var p users.UserPrimitives
	var email *string
	var memberships []byte
	err := row.Scan(&p.Id, &p.Name, &email, &p.SystemRole, &p.IsSystem, &p.Version, &p.AuthorId, &p.AuthorRole, &p.SubscriptionIds, &memberships)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	if email != nil {
		p.Email = *email
	}
	if err = json.Unmarshal(memberships, &p.Memberships); err != nil {
		return nil, fmt.Errorf("unmarshal user memberships: %w", err)
	}
	return users.FromPrimitives(p)
}
