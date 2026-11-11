package infrastructure_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared/infrastructure"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

type integrationIdentityProviderEncrypter struct{}

func (integrationIdentityProviderEncrypter) Encrypt(value string) (string, error) {
	return "encrypted:" + value, nil
}
func (integrationIdentityProviderEncrypter) Decrypt(value string) (string, error) {
	const prefix = "encrypted:"
	if len(value) < len(prefix) || value[:len(prefix)] != prefix {
		return "", fmt.Errorf("invalid encrypted value")
	}
	return value[len(prefix):], nil
}
func (integrationIdentityProviderEncrypter) Seal(value []byte) ([]byte, error) { return value, nil }
func (integrationIdentityProviderEncrypter) Open(value []byte) ([]byte, error) { return value, nil }

func TestPostgresProtocolIdentityProviderRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	id, referenceName := uuid.NewString(), "protocol-idp-"+uuid.NewString()
	_, err := pool.Exec(ctx, `INSERT INTO backoffice.identity_providers (id,reference_name,name,issuer,client_id,client_secret,authorization_endpoint,token_endpoint,jwks_uri,scopes,backoffice_redirect_uri) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		id, referenceName, "Integration Provider", "https://issuer.example/"+id, "client-id", "encrypted:client-secret", "https://issuer.example/authorize", "https://issuer.example/token", "https://issuer.example/jwks", []string{"openid", "email"}, "https://app.example/callback")
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM backoffice.identity_providers WHERE id=$1`, id) })

	repo := infrastructure.NewPostgresIdentityProviderRepository(pool, integrationIdentityProviderEncrypter{})
	config, err := repo.FindByReferenceName(ctx, referenceName)
	require.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, id, config.ID)
	assert.Equal(t, "client-secret", config.ClientSecret)
	assert.Equal(t, []string{"openid", "email"}, config.Scopes)
	assert.Equal(t, "https://app.example/callback", config.BackofficeRedirectUri)

	missing, err := repo.FindByReferenceName(ctx, "missing-"+referenceName)
	require.NoError(t, err)
	assert.Nil(t, missing)
}
