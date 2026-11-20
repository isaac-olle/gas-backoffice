package infrastructure

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	sharedInfrastructure "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSessionRepository struct {
	postgres  *sharedInfrastructure.PostgresRepository[*auth.Session]
	encrypter persistance.Encrypter
}

func NewPostgresSessionRepository(db *pgxpool.Pool, encrypter persistance.Encrypter) *PostgresSessionRepository {
	return &PostgresSessionRepository{
		postgres:  sharedInfrastructure.NewPostgresRepository[*auth.Session](db),
		encrypter: encrypter,
	}
}

func (r *PostgresSessionRepository) SaveSession(ctx context.Context, session *auth.Session) error {
	p := session.ToPrimitives()
	var encryptedCSRFToken *string
	if p.CSRFToken != "" {
		encrypted, err := r.encrypter.Encrypt(p.CSRFToken)
		if err != nil {
			return fmt.Errorf("encrypt backoffice session CSRF token: %w", err)
		}
		encryptedCSRFToken = &encrypted
	}
	_, err := r.postgres.Exec(ctx, `
		INSERT INTO backoffice.sessions (id, user_id, user_role, csrf_token)
		VALUES ($1, $2, $3, $4)
	`, p.SessionId, p.UserId, p.Role, encryptedCSRFToken)
	if err != nil {
		return fmt.Errorf("save backoffice session: %w", err)
	}
	return nil
}

func (r *PostgresSessionRepository) IsActive(ctx context.Context, id ids.SessionId) (bool, error) {
	var active bool
	err := r.postgres.QueryRow(ctx, `
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
	_, err := r.postgres.Exec(ctx, `
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
