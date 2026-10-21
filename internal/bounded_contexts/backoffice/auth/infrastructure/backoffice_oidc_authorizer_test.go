package infrastructure_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	backofficeClientId = "backoffice-client-123"
	backofficeCallback = "http://localhost:8090/v1/auth/callback"
	idpRedirectURL     = "https://accounts.google.com/o/oauth2/v2/auth?client_id=google-123&state=abc"
)

// passthroughEncrypter no xifra, però imita el sobre de l'AESEncrypter real:
// retorna base64 estàndard, que és el que el transcodificador d'state espera.
type passthroughEncrypter struct{ err error }

func (e passthroughEncrypter) Encrypt(plaintext string) (string, error) {
	if e.err != nil {
		return "", e.err
	}
	return base64.StdEncoding.EncodeToString([]byte(plaintext)), nil
}

func (e passthroughEncrypter) Decrypt(ciphertext string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (e passthroughEncrypter) Seal(plaintext []byte) ([]byte, error) {
	if e.err != nil {
		return nil, e.err
	}
	return plaintext, nil
}

func (e passthroughEncrypter) Open(sealed []byte) ([]byte, error) {
	return sealed, nil
}

// statePayloadFrom desfà el transport del state: URL-safe base64 -> el que hi
// hagi posat l'encrypter (aquí, el JSON en clar).
func statePayloadFrom(t *testing.T, state string) []byte {
	t.Helper()
	raw, err := base64.RawURLEncoding.DecodeString(state)
	require.NoError(t, err)
	return raw
}

func anIdp() auth.IdpConfig {
	return auth.IdpConfig{}
}

func aValidVerifier() string { return strings.Repeat("a", 43) }

type stubCodeVerifierGenerator struct{ value string }

func (g stubCodeVerifierGenerator) GenerateCodeVerifier() (vos.CodeVerifier, error) {
	return vos.NewCodeVerifierFromString(g.value)
}

func aLoginRequestWithVerifier(t *testing.T, verifier string) *auth.LoginRequest {
	t.Helper()
	idpReferenceName, err := vos.NewFromString("google")
	require.NoError(t, err)

	loginRequest, err := auth.NewLoginRequest(idpReferenceName, stubCodeVerifierGenerator{value: verifier})
	require.NoError(t, err)
	return loginRequest
}

// authorizeServerStub imita /oidc/authorize i captura la petició rebuda.
type authorizeServerStub struct {
	server   *httptest.Server
	received *http.Request
}

func newAuthorizeServerStub(t *testing.T, handler http.HandlerFunc) *authorizeServerStub {
	t.Helper()
	stub := &authorizeServerStub{}
	stub.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stub.received = r
		handler(w, r)
	}))
	t.Cleanup(stub.server.Close)
	return stub
}

func redirectingTo(location string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusFound)
	}
}

func anAuthorizerFor(endpoint string, encrypter passthroughEncrypter) *infrastructure.OidcAuthorizer {
	return infrastructure.NewOidcAuthorizer(encrypter, endpoint, backofficeClientId, backofficeCallback)
}

func TestAuthorize(t *testing.T) {
	t.Parallel()

	t.Run("returns the Location the authorization server redirects to", func(t *testing.T) {
		t.Parallel()
		stub := newAuthorizeServerStub(t, redirectingTo(idpRedirectURL))

		got, err := anAuthorizerFor(stub.server.URL, passthroughEncrypter{}).
			Authorize(context.Background(), anIdp(), aLoginRequestWithVerifier(t, aValidVerifier()))
		require.NoError(t, err)

		assert.Equal(t, idpRedirectURL, got, "no ha de seguir el redirect, només tornar-lo")
	})

	t.Run("sends the backoffice client credentials and the idp to federate to", func(t *testing.T) {
		t.Parallel()
		stub := newAuthorizeServerStub(t, redirectingTo(idpRedirectURL))

		_, err := anAuthorizerFor(stub.server.URL, passthroughEncrypter{}).
			Authorize(context.Background(), anIdp(), aLoginRequestWithVerifier(t, aValidVerifier()))
		require.NoError(t, err)

		require.NotNil(t, stub.received)
		params := stub.received.URL.Query()
		assert.Equal(t, http.MethodGet, stub.received.Method)
		assert.Equal(t, "code", params.Get("response_type"))
		assert.Equal(t, backofficeClientId, params.Get("client_id"))
		assert.Equal(t, backofficeCallback, params.Get("redirect_uri"))
		assert.Equal(t, "openid", params.Get("scope"))
		assert.Equal(t, "google", params.Get("idp_reference_name"))
	})

	t.Run("derives the S256 challenge from the login request verifier", func(t *testing.T) {
		t.Parallel()
		stub := newAuthorizeServerStub(t, redirectingTo(idpRedirectURL))

		_, err := anAuthorizerFor(stub.server.URL, passthroughEncrypter{}).
			Authorize(context.Background(), anIdp(), aLoginRequestWithVerifier(t, aValidVerifier()))
		require.NoError(t, err)

		params := stub.received.URL.Query()
		sum := sha256.Sum256([]byte(aValidVerifier()))
		assert.Equal(t, base64.RawURLEncoding.EncodeToString(sum[:]), params.Get("code_challenge"))
		assert.Equal(t, "S256", params.Get("code_challenge_method"))
	})

	t.Run("carries the verifier inside the state", func(t *testing.T) {
		t.Parallel()
		stub := newAuthorizeServerStub(t, redirectingTo(idpRedirectURL))

		_, err := anAuthorizerFor(stub.server.URL, passthroughEncrypter{}).
			Authorize(context.Background(), anIdp(), aLoginRequestWithVerifier(t, aValidVerifier()))
		require.NoError(t, err)

		var payload struct {
			CV  string `json:"cv"`
			Exp int64  `json:"exp"`
		}
		state := stub.received.URL.Query().Get("state")
		require.NoError(t, json.Unmarshal(statePayloadFrom(t, state), &payload))
		assert.Equal(t, aValidVerifier(), payload.CV)
		assert.NotZero(t, payload.Exp)
	})

	t.Run("fails when the authorization server does not redirect", func(t *testing.T) {
		t.Parallel()
		stub := newAuthorizeServerStub(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid_client"}`))
		})

		_, err := anAuthorizerFor(stub.server.URL, passthroughEncrypter{}).
			Authorize(context.Background(), anIdp(), aLoginRequestWithVerifier(t, aValidVerifier()))

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid_client")
	})

	t.Run("fails when the redirect carries no Location", func(t *testing.T) {
		t.Parallel()
		stub := newAuthorizeServerStub(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusFound)
		})

		_, err := anAuthorizerFor(stub.server.URL, passthroughEncrypter{}).
			Authorize(context.Background(), anIdp(), aLoginRequestWithVerifier(t, aValidVerifier()))

		require.Error(t, err)
		assert.Contains(t, err.Error(), "Location")
	})

	t.Run("fails when the authorization server is unreachable", func(t *testing.T) {
		t.Parallel()
		_, err := anAuthorizerFor("http://127.0.0.1:1/oidc/authorize", passthroughEncrypter{}).
			Authorize(context.Background(), anIdp(), aLoginRequestWithVerifier(t, aValidVerifier()))

		assert.Error(t, err)
	})

	t.Run("emits a url-safe state", func(t *testing.T) {
		t.Parallel()
		stub := newAuthorizeServerStub(t, redirectingTo(idpRedirectURL))

		_, err := anAuthorizerFor(stub.server.URL, passthroughEncrypter{}).
			Authorize(context.Background(), anIdp(), aLoginRequestWithVerifier(t, aValidVerifier()))
		require.NoError(t, err)

		state := stub.received.URL.Query().Get("state")
		assert.NotContains(t, state, "+", "un '+' en una query string es llegeix com un espai")
		assert.NotContains(t, state, "/")
		assert.NotContains(t, state, "=")
	})

	t.Run("fails when the state cannot be encrypted", func(t *testing.T) {
		t.Parallel()
		stub := newAuthorizeServerStub(t, redirectingTo(idpRedirectURL))

		_, err := anAuthorizerFor(stub.server.URL, passthroughEncrypter{err: fmt.Errorf("boom")}).
			Authorize(context.Background(), anIdp(), aLoginRequestWithVerifier(t, aValidVerifier()))

		assert.Error(t, err)
	})
}
