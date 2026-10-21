package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresIdentityProviderRepository struct {
	postgres *infrastructure.PostgresRepository[*identity_providers.IdentityProvider]
	pool     *pgxpool.Pool
}

func NewPostgresIdentityProviderRepository(db *pgxpool.Pool) *PostgresIdentityProviderRepository {
	return &PostgresIdentityProviderRepository{
		postgres: infrastructure.NewPostgresRepository[*identity_providers.IdentityProvider](db),
		pool:     db,
	}
}

func (r *PostgresIdentityProviderRepository) Save(ctx context.Context, idp *identity_providers.IdentityProvider) error {
	p := idp.ToPrimitives()
	_, err := r.pool.Exec(ctx, `
		INSERT INTO backoffice.identity_providers (
			id, reference_name, name, issuer,
			client_id, client_secret,
			authorization_endpoint, token_endpoint, userinfo_endpoint, jwks_uri,
			scopes,
			backoffice_redirect_uri,
			enabled
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`,
		p.Id, p.ReferenceName, p.Name, p.Issuer,
		p.ClientId, nullIfEmpty(p.ClientSecret),
		p.AuthorizationEndpoint, p.TokenEndpoint, nullIfEmpty(p.UserinfoEndpoint), p.JwksUri,
		p.Scopes,
		nullIfEmpty(p.BackofficeRedirectUri),
		p.Enabled,
	)
	if err != nil {
		return fmt.Errorf("save identity provider: %w", err)
	}
	return nil
}

func (r *PostgresIdentityProviderRepository) FindById(ctx context.Context, id ids.IdentityProviderId) (*identity_providers.IdentityProvider, error) {
	query := `
		SELECT id, reference_name, name, issuer, client_id, client_secret,
		       authorization_endpoint, token_endpoint, userinfo_endpoint, jwks_uri,
		       scopes, backoffice_redirect_uri, enabled
		FROM backoffice.identity_providers
		WHERE id = $1
	`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresIdentityProviderRepository) FindByReferenceName(ctx context.Context, referenceName string) (*identity_providers.IdentityProvider, error) {
	query := `
		SELECT id, reference_name, name, issuer, client_id, client_secret,
		       authorization_endpoint, token_endpoint, userinfo_endpoint, jwks_uri,
		       scopes, backoffice_redirect_uri, enabled
		FROM backoffice.identity_providers
		WHERE reference_name = $1
	`
	return r.postgres.Find(ctx, query, r.fromRow, referenceName)
}

func (r *PostgresIdentityProviderRepository) fromRow(row pgx.Row) (*identity_providers.IdentityProvider, error) {
	var p identity_providers.IdentityProviderPrimitives
	var clientSecret, userinfoEndpoint *string
	var backofficeRedirectUri *string
	err := row.Scan(
		&p.Id, &p.ReferenceName, &p.Name, &p.Issuer, &p.ClientId, &clientSecret,
		&p.AuthorizationEndpoint, &p.TokenEndpoint, &userinfoEndpoint, &p.JwksUri,
		&p.Scopes, &backofficeRedirectUri, &p.Enabled,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning identity provider: %w", err)
	}
	p.ClientSecret = deref(clientSecret)
	p.UserinfoEndpoint = deref(userinfoEndpoint)
	p.BackofficeRedirectUri = deref(backofficeRedirectUri)
	return identity_providers.FromPrimitives(p)
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
