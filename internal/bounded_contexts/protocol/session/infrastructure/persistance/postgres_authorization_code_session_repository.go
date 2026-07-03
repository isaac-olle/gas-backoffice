package persistance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAuthorizationCodeSessionRepository struct {
	postgres *infrastructure.PostgresRepository[*session.AuthorizationCodeSession]
}

func NewPostgresAuthorizationCodeSessionRepository(pool *pgxpool.Pool) *PostgresAuthorizationCodeSessionRepository {
	return &PostgresAuthorizationCodeSessionRepository{
		postgres: infrastructure.NewPostgresRepository[*session.AuthorizationCodeSession](pool),
	}
}

func (r *PostgresAuthorizationCodeSessionRepository) Save(ctx context.Context, sess *session.AuthorizationCodeSession) error {
	p := sess.ToPrimitives()

	pkceJSON, err := json.Marshal(p.Pkce)
	if err != nil {
		return fmt.Errorf("marshal pkce: %w", err)
	}

	responseTypesJSON, err := json.Marshal(p.ResponseTypes)
	if err != nil {
		return fmt.Errorf("marshal response types: %w", err)
	}

	query := `
		INSERT INTO protocol_authorization_code_sessions (
			session_id, client_id, redirect_uri, scope, state,
			response_types, code, pkce, session_state, flow, consumed,
			expires_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11, NOW() + interval '15 minutes')
		ON CONFLICT (session_id) DO UPDATE SET
			session_state = EXCLUDED.session_state,
			consumed = EXCLUDED.consumed,
			code = EXCLUDED.code
	`
	if err = r.postgres.Save(ctx, query,
		p.SessionID, p.ClientID, p.RedirectURI, p.Scope, p.State,
		responseTypesJSON, p.Code, pkceJSON, p.SessionState, p.Flow, p.Consumed,
	); err != nil {
		return fmt.Errorf("exec insert session: %w", err)
	}
	return nil
}

func (r *PostgresAuthorizationCodeSessionRepository) FindById(ctx context.Context, id ids.SessionId) (*session.AuthorizationCodeSession, error) {
	query := `
		SELECT session_id, client_id, redirect_uri, scope, state,
		       response_types, code, pkce, session_state, flow, consumed
		FROM protocol_authorization_code_sessions WHERE session_id = $1 AND deleted_at IS NULL
	`
	return r.postgres.Find(ctx, query, scanAuthorizationCodeSession, id.String())
}

func (r *PostgresAuthorizationCodeSessionRepository) FindByCode(ctx context.Context, code string) (*session.AuthorizationCodeSession, error) {
	query := `
		SELECT session_id, client_id, redirect_uri, scope, state,
		       response_types, code, pkce, session_state, flow, consumed
		FROM protocol_authorization_code_sessions WHERE code = $1 AND deleted_at IS NULL
	`
	return r.postgres.Find(ctx, query, scanAuthorizationCodeSession, code)
}

func (r *PostgresAuthorizationCodeSessionRepository) SoftDelete(ctx context.Context, id ids.SessionId) error {
	query := `UPDATE protocol_authorization_code_sessions SET deleted_at = NOW() WHERE session_id = $1`
	if err := r.postgres.Delete(ctx, query, id.String()); err != nil {
		return fmt.Errorf("soft delete session: %w", err)
	}
	return nil
}

func (r *PostgresAuthorizationCodeSessionRepository) HardDelete(ctx context.Context, id ids.SessionId) error {
	query := `DELETE FROM protocol_authorization_code_sessions WHERE session_id = $1`
	if err := r.postgres.Delete(ctx, query, id.String()); err != nil {
		return fmt.Errorf("hard delete session: %w", err)
	}
	return nil
}

func scanAuthorizationCodeSession(row pgx.Row) (*session.AuthorizationCodeSession, error) {
	var p session.AuthorizationCodeSessionPrimitives
	var responseTypesJSON, pkceJSON []byte

	err := row.Scan(
		&p.SessionID, &p.ClientID, &p.RedirectURI, &p.Scope, &p.State,
		&responseTypesJSON, &p.Code, &pkceJSON, &p.SessionState, &p.Flow, &p.Consumed,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan session: %w", err)
	}

	if err = json.Unmarshal(responseTypesJSON, &p.ResponseTypes); err != nil {
		return nil, fmt.Errorf("unmarshal response types: %w", err)
	}

	if pkceJSON != nil && string(pkceJSON) != "null" {
		p.Pkce = &session.PkcePrimitives{}
		if err = json.Unmarshal(pkceJSON, p.Pkce); err != nil {
			return nil, fmt.Errorf("unmarshal pkce: %w", err)
		}
	}

	return session.AuthorizationCodeSessionFromPrimitives(&p)
}
