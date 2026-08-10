package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSessionRepository struct {
	postgres *infrastructure.PostgresRepository[*auth.Session]
}

func NewPostgresSessionRepository(db *pgxpool.Pool) *PostgresSessionRepository {
	return &PostgresSessionRepository{
		postgres: infrastructure.NewPostgresRepository[*auth.Session](db),
	}
}

func (r *PostgresSessionRepository) Save(ctx context.Context, session *auth.Session) error {
	p := session.ToPrimitives()
	authorId, authorRole := authorPrimitives(ctx)
	query := `
		INSERT INTO backoffice.sessions (id, user_id, author_id, author_role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	return r.postgres.Save(ctx, query, p.Id, p.UserId, authorId, authorRole)
}

func authorPrimitives(ctx context.Context) (string, string) {
	author, ok := vo.AuthorFromContext[vo.Author](ctx)
	if !ok {
		return "", ""
	}
	p := author.ToPrimitives()
	return p.Id, p.Role
}

func (r *PostgresSessionRepository) SearchById(ctx context.Context, id ids.SessionId) (*auth.Session, error) {
	query := `SELECT id, user_id FROM backoffice.sessions WHERE id = $1`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresSessionRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*auth.Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresSessionRepository) SoftDelete(ctx context.Context, _ ids.SessionId) error {
	return fmt.Errorf("not implemented")
}

func (r *PostgresSessionRepository) HardDelete(ctx context.Context, id ids.SessionId) error {
	query := `DELETE FROM backoffice.sessions WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresSessionRepository) fromRow(row pgx.Row) (*auth.Session, error) {
	var p auth.SessionPrimitives
	err := row.Scan(&p.Id, &p.UserId)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning session: %w", err)
	}
	return auth.FromPrimitives(p)
}
