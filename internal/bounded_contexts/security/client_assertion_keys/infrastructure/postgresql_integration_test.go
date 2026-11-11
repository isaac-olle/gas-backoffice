package infrastructure_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/infrastructure"
	keyPairInfrastructure "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

const integrationPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArRj90RkbqGVRcp1xYi63
L+ZKaqMi3IVFhm6adSM+SNJsdiEkcMjrXu8iKlxN19kfJFN/fDlKwZuDgdKLq2/X
QqKpQHItF6kgNbuzfkNWbLkzebWn5sokC9BqgStH4Inob5Plg8T4v2Po7JIn8098
bGtxNPLObG/+heX7JQNZDyocWpQWD4G4BrvgimfKGH4azvBmLHR4dayCrii7GR2I
W1kY4MmhX84TNLCqFnzZdd0GswyGJHVnyrpUvEJFaHmr+97I8vDBvFPJ2w2dEXN1
+BC5dJHqmLjr6qV+6Dqq56Gn8DcestRxtiS20np4Z2cMSCgqtgp1l1zHSjMO/rGX
WQIDAQAB
-----END PUBLIC KEY-----
`

func TestPostgresClientAssertionKeyRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, rawSubscriptionID := integrationtest.BackofficeSubscription(t, pool)
	subscriptionID, err := ids.NewSubscriptionIdFromString(rawSubscriptionID)
	require.NoError(t, err)
	repo := infrastructure.NewPostgresClientAssertionKeyRepository(pool)
	clientID := ids.NewClientId()
	keyID := ids.NewClientAssertionKeyId()
	t.Cleanup(func() { _ = repo.HardDelete(ctx, keyID) })

	key, err := client_assertion_keys.NewClientAssertionKeyForSubscription(keyID, clientID, subscriptionID, "RS256", integrationPublicKey)
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, key))

	byID, err := repo.SearchById(ctx, keyID)
	require.NoError(t, err)
	require.NotNil(t, byID)
	assert.Equal(t, "RS256", byID.ToPrimitives().Algorithm)
	assert.Equal(t, []byte(integrationPublicKey), byID.ToPrimitives().PublicKey)

	byClient, err := repo.FindByClientId(ctx, clientID)
	require.NoError(t, err)
	require.NotNil(t, byClient)
	assert.Equal(t, keyID.String(), byClient.ToPrimitives().Id)

	require.NoError(t, repo.SoftDelete(ctx, keyID))
	missing, err := repo.SearchById(ctx, keyID)
	require.NoError(t, err)
	assert.Nil(t, missing)
}

func TestJwtAssertionVerifierWithPostgresKeyRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	_, rawSubscriptionID := integrationtest.BackofficeSubscription(t, pool)
	subscriptionID, err := ids.NewSubscriptionIdFromString(rawSubscriptionID)
	require.NoError(t, err)
	clientID := ids.NewClientId()
	keyID := ids.NewClientAssertionKeyId()
	privatePEM, publicPEM, err := keyPairInfrastructure.NewPEMKeyGenerator().Generate("ES256")
	require.NoError(t, err)
	key, err := client_assertion_keys.NewClientAssertionKeyForSubscription(keyID, clientID, subscriptionID, "ES256", publicPEM)
	require.NoError(t, err)
	repo := infrastructure.NewPostgresClientAssertionKeyRepository(pool)
	t.Cleanup(func() { _ = repo.HardDelete(ctx, keyID) })
	require.NoError(t, repo.Save(ctx, key))

	block, _ := pem.Decode([]byte(privatePEM))
	require.NotNil(t, block)
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	require.NoError(t, err)
	privateKey, ok := parsed.(*ecdsa.PrivateKey)
	require.True(t, ok)
	makeAssertion := func(key *ecdsa.PrivateKey) string {
		t.Helper()
		claims := jwt.MapClaims{
			"iss": clientID.String(), "sub": clientID.String(), "aud": "https://api.example/token",
			"jti": "assertion-1", "exp": time.Now().Add(time.Minute).Unix(),
		}
		assertion, signErr := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(key)
		require.NoError(t, signErr)
		return assertion
	}
	verifier := infrastructure.NewJwtAssertionVerifier(domain_services.NewClientAssertionKeyFinder(repo))
	verified, err := verifier.Verify(ctx, makeAssertion(privateKey))
	require.NoError(t, err)
	assert.Equal(t, clientID.String(), verified["sub"])

	wrongPrivatePEM, _, err := keyPairInfrastructure.NewPEMKeyGenerator().Generate("ES256")
	require.NoError(t, err)
	wrongBlock, _ := pem.Decode([]byte(wrongPrivatePEM))
	require.NotNil(t, wrongBlock)
	wrongParsed, err := x509.ParsePKCS8PrivateKey(wrongBlock.Bytes)
	require.NoError(t, err)
	wrongPrivateKey, ok := wrongParsed.(*ecdsa.PrivateKey)
	require.True(t, ok)
	_, err = verifier.Verify(ctx, makeAssertion(wrongPrivateKey))
	assert.Error(t, err)
}
