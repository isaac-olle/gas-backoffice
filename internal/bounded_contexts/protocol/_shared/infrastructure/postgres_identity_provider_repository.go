package infrastructure

import (
	"context"
	"errors"
	"fmt"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresIdentityProviderRepository struct {
	pool      *pgxpool.Pool
	encrypter persistance.Encrypter
}

func NewPostgresIdentityProviderRepository(pool *pgxpool.Pool, encrypter persistance.Encrypter) *PostgresIdentityProviderRepository {
	return &PostgresIdentityProviderRepository{pool: pool, encrypter: encrypter}
}

func (r *PostgresIdentityProviderRepository) FindByReferenceName(ctx context.Context, referenceName string) (*protocolShared.IdentityProviderConfig, error) {
	var (
		id, name, issuer, clientID, authorizationEndpoint, tokenEndpoint, jwksUri string
		clientSecret, backofficeRedirectUri                                       *string
		scopes                                                                    []string
	)
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, issuer, client_id, client_secret, authorization_endpoint, token_endpoint, jwks_uri, scopes, backoffice_redirect_uri
		FROM backoffice.identity_providers
		WHERE reference_name = $1 AND enabled = TRUE
	`, referenceName).Scan(
		&id, &name, &issuer, &clientID, &clientSecret,
		&authorizationEndpoint, &tokenEndpoint, &jwksUri, &scopes, &backofficeRedirectUri,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying identity provider: %w", err)
	}

	decryptedSecret := ""
	if clientSecret != nil && *clientSecret != "" {
		if decryptedSecret, err = r.encrypter.Decrypt(*clientSecret); err != nil {
			return nil, fmt.Errorf("error decrypting identity provider client secret: %w", err)
		}
	}

	redirectUri := ""
	if backofficeRedirectUri != nil {
		redirectUri = *backofficeRedirectUri
	}

	return protocolShared.NewIdentityProviderConfig(
		id, name, issuer, clientID, decryptedSecret,
		authorizationEndpoint, tokenEndpoint, jwksUri, scopes, redirectUri,
	), nil
}
