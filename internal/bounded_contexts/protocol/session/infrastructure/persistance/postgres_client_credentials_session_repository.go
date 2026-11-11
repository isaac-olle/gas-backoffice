package persistance

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedInfrastructure "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClientCredentialsSessionRepository struct {
	postgres *sharedInfrastructure.PostgresRepository[*session.ClientCredentialsSession]
}

func NewPostgresClientCredentialsSessionRepository(pool *pgxpool.Pool) *PostgresClientCredentialsSessionRepository {
	return &PostgresClientCredentialsSessionRepository{postgres: sharedInfrastructure.NewPostgresRepository[*session.ClientCredentialsSession](pool)}
}

func (r *PostgresClientCredentialsSessionRepository) Save(ctx context.Context, sess *session.ClientCredentialsSession) error {
	p := sess.ToPrimitives()

	_, err := r.postgres.Exec(ctx, `
		INSERT INTO protocol.client_credentials_sessions (
			id, client_id, user_id, scope, issued_at, protected_resource_id
		) VALUES ($1,$2,$3,$4,$5,$6)
	`, p.ID, p.ClientID, p.UserId, p.Scope, p.IssuedAt, p.ProtectedResourceId)
	if err != nil {
		return fmt.Errorf("exec insert client credentials session: %w", err)
	}
	return nil
}

func (r *PostgresClientCredentialsSessionRepository) FindById(ctx context.Context, id ids.SessionId) (*session.ClientCredentialsSession, error) {
	row := r.postgres.QueryRow(ctx, `
		SELECT id, client_id, user_id, scope, issued_at, protected_resource_id
		FROM protocol.client_credentials_sessions WHERE id = $1
	`, id.String())

	p := session.ClientCredentialsSessionPrimitives{BaseSessionPrimitives: &session.BaseSessionPrimitives{}}
	if err := row.Scan(&p.ID, &p.ClientID, &p.UserId, &p.Scope, &p.IssuedAt, &p.ProtectedResourceId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan client credentials session: %w", err)
	}
	return session.ClientCredentialsSessionFromPrimitives(&p)
}
