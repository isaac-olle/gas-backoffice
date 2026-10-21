package persistance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAuthorizationCodeSessionRepository struct {
	postgres  *infrastructure.PostgresRepository[*session.AuthorizationCodeSession]
	encrypter persistance.Encrypter
}

func NewPostgresAuthorizationCodeSessionRepository(pool *pgxpool.Pool, encrypter persistance.Encrypter) *PostgresAuthorizationCodeSessionRepository {
	return &PostgresAuthorizationCodeSessionRepository{
		postgres:  infrastructure.NewPostgresRepository[*session.AuthorizationCodeSession](pool),
		encrypter: encrypter,
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

	claimsJSON, err := json.Marshal(p.Identity.Claims)
	if err != nil {
		return fmt.Errorf("marshal identity claims: %w", err)
	}

	encryptedIDToken, err := r.encrypter.Encrypt(p.Identity.IDToken)
	if err != nil {
		return fmt.Errorf("encrypt id_token: %w", err)
	}

	encryptedAccessToken, err := r.encrypter.Encrypt(p.Identity.AccessToken)
	if err != nil {
		return fmt.Errorf("encrypt access_token: %w", err)
	}

	var encryptedRefreshToken *string
	if p.Identity.RefreshToken != nil && *p.Identity.RefreshToken != "" {
		encrypted, err := r.encrypter.Encrypt(*p.Identity.RefreshToken)
		if err != nil {
			return fmt.Errorf("encrypt refresh_token: %w", err)
		}
		encryptedRefreshToken = &encrypted
	}

	query := `
		INSERT INTO protocol.authorization_code_sessions (
			session_id, client_id, client_owner_user_id, redirect_uri, scope, state,
			response_types, code, pkce, session_state, grant_type, consumed,
			protected_resource_id, idp_id, id_token, access_token, refresh_token,
			claims, identity_expires_at, expires_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18::jsonb,$19, NOW() + interval '15 minutes')
		ON CONFLICT (session_id) DO UPDATE SET
			session_state = EXCLUDED.session_state,
			consumed      = EXCLUDED.consumed,
			code          = EXCLUDED.code
	`
	if err = r.postgres.Save(ctx, query,
		p.SessionID, p.ClientID, p.ClientOwnerUserId, p.RedirectURI, p.Scope, p.State,
		responseTypesJSON, p.Code, pkceJSON, p.SessionState, p.GrantType, p.Consumed,
		p.ProtectedResourceId, p.Identity.IdpId, encryptedIDToken, encryptedAccessToken, encryptedRefreshToken,
		claimsJSON, p.Identity.ExpiresAt,
	); err != nil {
		return fmt.Errorf("exec insert session: %w", err)
	}
	return nil
}

func (r *PostgresAuthorizationCodeSessionRepository) FindById(ctx context.Context, id ids.SessionId) (*session.AuthorizationCodeSession, error) {
	query := `
		SELECT session_id, client_id, client_owner_user_id, redirect_uri, scope, state,
		       response_types, code, pkce, session_state, grant_type, consumed, protected_resource_id,
		       idp_id, id_token, access_token, refresh_token, claims, identity_expires_at
		FROM protocol.authorization_code_sessions WHERE session_id = $1 AND deleted_at IS NULL
	`
	return r.postgres.Find(ctx, query, r.scanAuthorizationCodeSession, id.String())
}

func (r *PostgresAuthorizationCodeSessionRepository) FindByCode(ctx context.Context, code string) (*session.AuthorizationCodeSession, error) {
	query := `
		SELECT session_id, client_id, client_owner_user_id, redirect_uri, scope, state,
		       response_types, code, pkce, session_state, grant_type, consumed, protected_resource_id,
		       idp_id, id_token, access_token, refresh_token, claims, identity_expires_at
		FROM protocol.authorization_code_sessions WHERE code = $1 AND deleted_at IS NULL
	`
	return r.postgres.Find(ctx, query, r.scanAuthorizationCodeSession, code)
}

func (r *PostgresAuthorizationCodeSessionRepository) SoftDelete(ctx context.Context, id ids.SessionId) error {
	query := `UPDATE protocol.authorization_code_sessions SET deleted_at = NOW() WHERE session_id = $1`
	if err := r.postgres.Delete(ctx, query, id.String()); err != nil {
		return fmt.Errorf("soft delete session: %w", err)
	}
	return nil
}

func (r *PostgresAuthorizationCodeSessionRepository) HardDelete(ctx context.Context, id ids.SessionId) error {
	query := `DELETE FROM protocol.authorization_code_sessions WHERE session_id = $1`
	if err := r.postgres.Delete(ctx, query, id.String()); err != nil {
		return fmt.Errorf("hard delete session: %w", err)
	}
	return nil
}

func (r *PostgresAuthorizationCodeSessionRepository) scanAuthorizationCodeSession(row pgx.Row) (*session.AuthorizationCodeSession, error) {
	p := session.AuthorizationCodeSessionPrimitives{BaseSessionPrimitives: &session.BaseSessionPrimitives{}}
	var responseTypesJSON, pkceJSON, claimsJSON []byte
	var encryptedIDToken, encryptedAccessToken string
	var encryptedRefreshToken *string

	err := row.Scan(
		&p.SessionID, &p.ClientID, &p.ClientOwnerUserId, &p.RedirectURI, &p.Scope, &p.State,
		&responseTypesJSON, &p.Code, &pkceJSON, &p.SessionState, &p.GrantType, &p.Consumed,
		&p.ProtectedResourceId,
		&p.Identity.IdpId, &encryptedIDToken, &encryptedAccessToken, &encryptedRefreshToken,
		&claimsJSON, &p.Identity.ExpiresAt,
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

	if claimsJSON != nil {
		if err = json.Unmarshal(claimsJSON, &p.Identity.Claims); err != nil {
			return nil, fmt.Errorf("unmarshal identity claims: %w", err)
		}
	}

	if p.Identity.IDToken, err = r.encrypter.Decrypt(encryptedIDToken); err != nil {
		return nil, fmt.Errorf("decrypt id_token: %w", err)
	}
	if p.Identity.AccessToken, err = r.encrypter.Decrypt(encryptedAccessToken); err != nil {
		return nil, fmt.Errorf("decrypt access_token: %w", err)
	}
	if encryptedRefreshToken != nil && *encryptedRefreshToken != "" {
		decrypted, err := r.encrypter.Decrypt(*encryptedRefreshToken)
		if err != nil {
			return nil, fmt.Errorf("decrypt refresh_token: %w", err)
		}
		p.Identity.RefreshToken = &decrypted
	}

	return session.AuthorizationCodeSessionFromPrimitives(&p)
}
