package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
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

func (r *PostgresUserRepository) Save(ctx context.Context, user *users.User) error {
	p := user.ToPrimitives()
	authorId, authorRole := authorPrimitives(ctx)
	query := `
		INSERT INTO backoffice.users (id, name, admin_id, client_ids, protected_resource_ids, is_system, author_id, author_role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			client_ids = EXCLUDED.client_ids,
			protected_resource_ids = EXCLUDED.protected_resource_ids
	`
	return r.postgres.Save(ctx, query, p.Id, p.Name, p.AdminId, p.ClientIds, p.ProtectedResourceIds, p.IsSystem, authorId, authorRole)
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
	query := `SELECT id, name, admin_id, client_ids, protected_resource_ids, is_system FROM backoffice.users WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresUserRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*users.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresUserRepository) SoftDelete(ctx context.Context, id ids.UserId) error {
	query := `UPDATE backoffice.users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresUserRepository) HardDelete(ctx context.Context, id ids.UserId) error {
	query := `DELETE FROM backoffice.users WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresUserRepository) fromRow(row pgx.Row) (*users.User, error) {
	var p users.UserPrimitives
	err := row.Scan(&p.Id, &p.Name, &p.AdminId, &p.ClientIds, &p.ProtectedResourceIds, &p.IsSystem)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	return users.FromPrimitives(p)
}
