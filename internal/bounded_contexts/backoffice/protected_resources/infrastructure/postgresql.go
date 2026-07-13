package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProtectedResourceRepository struct {
	postgres *infrastructure.PostgresRepository[*protected_resources.ProtectedResource]
}

func NewPostgresProtectedResourceRepository(db *pgxpool.Pool) *PostgresProtectedResourceRepository {
	return &PostgresProtectedResourceRepository{
		postgres: infrastructure.NewPostgresRepository[*protected_resources.ProtectedResource](db),
	}
}

func (r *PostgresProtectedResourceRepository) Save(ctx context.Context, pr *protected_resources.ProtectedResource) error {
	p := pr.ToPrimitives()
	authorId, authorRole := actorPrimitives(ctx)
	query := `
		INSERT INTO protected_resources (id, name, user_id, admin_id, scope_ids, is_system, author_id, author_role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			scope_ids = EXCLUDED.scope_ids
	`
	return r.postgres.Save(ctx, query, p.Id, p.Name, p.UserId, p.AdminId, p.ScopeIds, p.IsSystem, authorId, authorRole)
}

func actorPrimitives(ctx context.Context) (string, string) {
	actor, ok := vo.ActorFromContext[vo.Actor](ctx)
	if !ok {
		return "", ""
	}
	p := actor.ToPrimitives()
	return p.Id, p.Role
}

func (r *PostgresProtectedResourceRepository) SearchById(ctx context.Context, id ids.ProtectedResourceId) (*protected_resources.ProtectedResource, error) {
	query := `SELECT id, name, user_id, admin_id, scope_ids, is_system FROM protected_resources WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresProtectedResourceRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*protected_resources.ProtectedResource, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresProtectedResourceRepository) SoftDelete(ctx context.Context, id ids.ProtectedResourceId) error {
	query := `UPDATE protected_resources SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresProtectedResourceRepository) HardDelete(ctx context.Context, id ids.ProtectedResourceId) error {
	query := `DELETE FROM protected_resources WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresProtectedResourceRepository) fromRow(row pgx.Row) (*protected_resources.ProtectedResource, error) {
	var p protected_resources.ProtectedResourcePrimitives
	err := row.Scan(&p.Id, &p.Name, &p.UserId, &p.AdminId, &p.ScopeIds, &p.IsSystem)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning protected resource: %w", err)
	}
	return protected_resources.FromPrimitives(p)
}
