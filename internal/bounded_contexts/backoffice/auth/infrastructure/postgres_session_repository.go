package infrastructure

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedInfrastructure "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSessionRepository struct {
	postgres *sharedInfrastructure.PostgresRepository[*auth.Session]
}

func NewPostgresSessionRepository(db *pgxpool.Pool) *PostgresSessionRepository {
	return &PostgresSessionRepository{postgres: sharedInfrastructure.NewPostgresRepository[*auth.Session](db)}
}

func (r *PostgresSessionRepository) SaveSession(ctx context.Context, session *auth.Session) error {
	p := session.ToPrimitives()
	_, err := r.postgres.DB().Exec(ctx, `
		INSERT INTO backoffice.sessions (id, user_id, user_role)
		VALUES ($1, $2, $3)
	`, p.SessionId, p.UserId, p.Role)
	if err != nil {
		return fmt.Errorf("save backoffice session: %w", err)
	}
	return nil
}

func (r *PostgresSessionRepository) IsActive(ctx context.Context, id ids.SessionId) (bool, error) {
	var active bool
	err := r.postgres.DB().QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM backoffice.sessions
			WHERE id = $1 AND logged_out_at IS NULL
		)
	`, id.String()).Scan(&active)
	if err != nil {
		return false, fmt.Errorf("check backoffice session: %w", err)
	}
	return active, nil
}

func (r *PostgresSessionRepository) Logout(ctx context.Context, id ids.SessionId) error {
	_, err := r.postgres.DB().Exec(ctx, `
		UPDATE backoffice.sessions
		SET logged_out_at = NOW()
		WHERE id = $1 AND logged_out_at IS NULL
	`, id.String())
	if err != nil {
		return fmt.Errorf("logout backoffice session: %w", err)
	}
	return nil
}

var _ auth.SessionRepository = (*PostgresSessionRepository)(nil)
