package persistance

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClientCredentialsSessionRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresClientCredentialsSessionRepository(pool *pgxpool.Pool) *PostgresClientCredentialsSessionRepository {
	return &PostgresClientCredentialsSessionRepository{pool: pool}
}

func (r *PostgresClientCredentialsSessionRepository) Save(ctx context.Context, sess *session.ClientCredentialsSession) error {
	p := sess.ToPrimitives()

	_, err := r.pool.Exec(ctx, `
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
	row := r.pool.QueryRow(ctx, `
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
