package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresIdentityProviderRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresIdentityProviderRepository(pool *pgxpool.Pool) *PostgresIdentityProviderRepository {
	return &PostgresIdentityProviderRepository{pool: pool}
}

func (r *PostgresIdentityProviderRepository) FindById(ctx context.Context, idpID string) (*session.IdentityProvider, error) {
	var (
		id, name, issuer, clientID  string
		authorizationEndpoint       string
		tokenEndpoint               string
		jwksURI                     string
		scopes                      []string
		clientSecret, userinfoEndpt *string
	)
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, issuer, client_id, client_secret, authorization_endpoint,
		       token_endpoint, userinfo_endpoint, jwks_uri, scopes
		FROM backoffice.oauth_identity_providers
		WHERE id = $1 AND enabled = TRUE
	`, idpID).Scan(
		&id, &name, &issuer, &clientID, &clientSecret, &authorizationEndpoint,
		&tokenEndpoint, &userinfoEndpt, &jwksURI, &scopes,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying identity provider: %w", err)
	}

	return session.NewIdentityProvider(
		id, name, issuer, clientID, deref(clientSecret),
		authorizationEndpoint, tokenEndpoint, deref(userinfoEndpt), jwksURI, scopes,
	), nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
