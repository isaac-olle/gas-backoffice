package infrastructure_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	sharedJwt "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/private_keys"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/public_keys"
)

const (
	gasIssuer   = "backoffice-protected-resource"
	gasClientId = "backoffice-client"
	gasKid      = "gas-kid"
)

type gasSigner struct {
	sign func(header, claims map[string]any) (string, error)
}

// aGasJwks serves a JWKS with a single RS256 key and returns its URL plus a
// signer for tokens with that key.
func aGasJwks(t *testing.T) (string, gasSigner) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privDer, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	signingKey, err := private_keys.NewRSAPrivateKeyFromPEM(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDer}), "RS256")
	require.NoError(t, err)

	pubDer, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	require.NoError(t, err)
	verifyingKey, err := public_keys.NewRSAPublicKeyFromPEM(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDer}), "RS256")
	require.NoError(t, err)
	jwk, err := verifyingKey.ToJWK()
	require.NoError(t, err)
	jwk["kid"] = gasKid

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]any{jwk}})
	}))
	t.Cleanup(server.Close)

	return server.URL, gasSigner{sign: sharedJwt.NewTokenSigner(signingKey).Sign}
}

func validClaims() map[string]any {
	return map[string]any{
		"iss":   gasIssuer,
		"aud":   gasClientId,
		"sub":   "user-sub",
		"email": "jane@example.com",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
}

func TestGasIdTokenVerifier(t *testing.T) {
	t.Parallel()
	header := map[string]any{"alg": "RS256", "typ": "JWT", "kid": gasKid}

	tests := []struct {
		name    string
		mutate  func(header, claims map[string]any)
		wantErr bool
	}{
		{name: "valid_token_returns_claims"},
		{name: "unknown_kid_errors", mutate: func(h, _ map[string]any) { h["kid"] = "other" }, wantErr: true},
		{name: "wrong_issuer_errors", mutate: func(_, c map[string]any) { c["iss"] = "https://accounts.google.com" }, wantErr: true},
		{name: "wrong_audience_errors", mutate: func(_, c map[string]any) { c["aud"] = "someone-else" }, wantErr: true},
		{name: "expired_token_errors", mutate: func(_, c map[string]any) { c["exp"] = time.Now().Add(-time.Minute).Unix() }, wantErr: true},
		{name: "missing_exp_errors", mutate: func(_, c map[string]any) { delete(c, "exp") }, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			jwksUri, signer := aGasJwks(t)
			h, c := map[string]any{}, validClaims()
			for k, v := range header {
				h[k] = v
			}
			if tt.mutate != nil {
				tt.mutate(h, c)
			}
			token, err := signer.sign(h, c)
			require.NoError(t, err)

			claims, err := infrastructure.NewGasIdTokenVerifier(jwksUri, gasIssuer, gasClientId).Verify(context.Background(), token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "jane@example.com", claims["email"])
		})
	}

	t.Run("tampered_signature_errors", func(t *testing.T) {
		t.Parallel()
		jwksUri, _ := aGasJwks(t)
		_, otherSigner := aGasJwks(t)
		token, err := otherSigner.sign(header, validClaims())
		require.NoError(t, err)

		_, err = infrastructure.NewGasIdTokenVerifier(jwksUri, gasIssuer, gasClientId).Verify(context.Background(), token)
		assert.Error(t, err)
	})
}
