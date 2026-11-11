package infrastructure_test

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure/helpers"
	keyPairInfrastructure "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/key_factory"
)

func TestOidcIdentityProviderClientExchangeCode(t *testing.T) {
	t.Parallel()
	privatePEM, publicPEM, err := keyPairInfrastructure.NewPEMKeyGenerator().Generate("RS256")
	require.NoError(t, err)
	block, _ := pem.Decode([]byte(privatePEM))
	require.NotNil(t, block)
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	require.NoError(t, err)
	privateKey, ok := parsed.(*rsa.PrivateKey)
	require.True(t, ok)

	const kid = "integration-key"
	verifyingKey, err := key_factory.NewVerifyingKeyFromPEM([]byte(publicPEM), "RS256")
	require.NoError(t, err)
	jwk, err := verifyingKey.ToJWK()
	require.NoError(t, err)
	jwk["kid"] = kid
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{jwk}})
	}))
	defer jwksServer.Close()

	const issuer = "https://issuer.example"
	const clientID = "oidc-client"
	idTokenClaims := jwt.MapClaims{
		"iss": issuer, "aud": clientID, "sub": "provider-user", "email": "user@example.test",
		"picture": "https://images.example/user.png", "exp": time.Now().Add(time.Minute).Unix(),
	}
	idTokenJWT := jwt.NewWithClaims(jwt.SigningMethodRS256, idTokenClaims)
	idTokenJWT.Header["kid"] = kid
	idToken, err := idTokenJWT.SignedString(privateKey)
	require.NoError(t, err)

	var requestForm url.Values
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		if parseErr := r.ParseForm(); parseErr == nil {
			requestForm = r.PostForm
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id_token": idToken, "access_token": "access-token", "refresh_token": "refresh-token", "expires_in": 300,
		})
	}))
	defer tokenServer.Close()

	idp := protocolShared.NewIdentityProviderConfig("idp-id", "Google", issuer, clientID, "client-secret", "", tokenServer.URL, jwksServer.URL, []string{"openid", "email"}, "https://app.example/callback")
	client := infrastructure.NewOidcIdentityProviderClient(helpers.NewOidcTokenVerifier())
	result, err := client.ExchangeCode(context.Background(), idp, "auth-code")
	require.NoError(t, err)
	assert.Equal(t, "auth-code", requestForm.Get("code"))
	assert.Equal(t, clientID, requestForm.Get("client_id"))
	assert.Equal(t, "client-secret", requestForm.Get("client_secret"))
	assert.Equal(t, "https://app.example/callback", requestForm.Get("redirect_uri"))
	assert.Equal(t, "authorization_code", requestForm.Get("grant_type"))
	assert.Equal(t, "access-token", result.Tokens.AccessToken)
	assert.Equal(t, "refresh-token", result.Tokens.RefreshToken)
	assert.Equal(t, idToken, result.Tokens.IDToken)
	assert.Equal(t, "provider-user", result.Claims.Sub)
	assert.Equal(t, "user@example.test", result.Claims.Email)
	assert.Equal(t, "https://images.example/user.png", result.Claims.Picture)
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), result.Tokens.ExpiresAt, 2*time.Second)

}
