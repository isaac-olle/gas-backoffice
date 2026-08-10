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
}

func NewPostgresIdentityProviderRepository(db *pgxpool.Pool) *PostgresIdentityProviderRepository {
	return &PostgresIdentityProviderRepository{
		postgres: infrastructure.NewPostgresRepository[*identity_providers.IdentityProvider](db),
	}
}

func (r *PostgresIdentityProviderRepository) FindById(ctx context.Context, id ids.IdentityProviderId) (*identity_providers.IdentityProvider, error) {
	query := `
		SELECT id, name, issuer, client_id, client_secret, authorization_endpoint,
		       token_endpoint, userinfo_endpoint, jwks_uri, scopes, enabled
		FROM backoffice.oauth_identity_providers
		WHERE id = $1
	`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresIdentityProviderRepository) fromRow(row pgx.Row) (*identity_providers.IdentityProvider, error) {
	var p identity_providers.IdentityProviderPrimitives
	var clientSecret, userinfoEndpoint *string
	err := row.Scan(
		&p.Id, &p.Name, &p.Issuer, &p.ClientId, &clientSecret, &p.AuthorizationEndpoint,
		&p.TokenEndpoint, &userinfoEndpoint, &p.JwksUri, &p.Scopes, &p.Enabled,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning identity provider: %w", err)
	}
	if clientSecret != nil {
		p.ClientSecret = *clientSecret
	}
	if userinfoEndpoint != nil {
		p.UserinfoEndpoint = *userinfoEndpoint
	}
	return identity_providers.FromPrimitives(p)
}
