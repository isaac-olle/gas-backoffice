//go:build e2e

package e2e_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminEndpoints(t *testing.T) {
	auth := withBearer(testAdminToken)

	userID := newID()
	prID := newID()
	scopeID := newID()
	clientID := newID()
	assertionKeyID := newID()
	secretID := newID()

	_, clientPubKeyPEM, err := generateRSAKeyPair()
	require.NoError(t, err)

	t.Run("GetAdmin", func(t *testing.T) {
		resp := do(http.MethodGet, "/v1/backoffice/admins", nil, auth)
		drain(resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("CreateUser", func(t *testing.T) {
		resp := do(http.MethodPut,
			fmt.Sprintf("/v1/backoffice/admins/users/%s", userID),
			map[string]any{"name": "e2e-user"},
			auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("GetUser", func(t *testing.T) {
		resp := do(http.MethodGet,
			fmt.Sprintf("/v1/backoffice/admins/users/%s", userID),
			nil, auth,
		)
		var body map[string]any
		decodeBody(resp, &body)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, userID, body["id"])
	})

	t.Run("ChangeUserName", func(t *testing.T) {
		resp := do(http.MethodPatch,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/name", userID),
			map[string]any{"name": "e2e-user-renamed"},
			auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("CreateProtectedResource", func(t *testing.T) {
		resp := do(http.MethodPut,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/protected-resources/%s", userID, prID),
			map[string]any{"name": "e2e-pr", "scopes": []any{}},
			auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("GetProtectedResource", func(t *testing.T) {
		resp := do(http.MethodGet,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/protected-resources/%s", userID, prID),
			nil, auth,
		)
		var body map[string]any
		decodeBody(resp, &body)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, prID, body["id"])
	})

	t.Run("UpdateProtectedResourceName", func(t *testing.T) {
		resp := do(http.MethodPatch,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/protected-resources/%s/name", userID, prID),
			map[string]any{"name": "e2e-pr-renamed"},
			auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("CreateScope", func(t *testing.T) {
		resp := do(http.MethodPut,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/protected-resources/%s/scopes/%s", userID, prID, scopeID),
			map[string]any{"name": "e2e:read"},
			auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("GetScope", func(t *testing.T) {
		resp := do(http.MethodGet,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/protected-resources/%s/scopes/%s", userID, prID, scopeID),
			nil, auth,
		)
		var body map[string]any
		decodeBody(resp, &body)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, scopeID, body["id"])
	})

	t.Run("CreateClient", func(t *testing.T) {
		resp := do(http.MethodPut,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/clients/%s", userID, clientID),
			map[string]any{
				"name":             "e2e-client",
				"assertion_key_id": assertionKeyID,
				"public_key":       clientPubKeyPEM,
				"secret_id":        secretID,
				"plaintext_secret": "super-secret-passphrase-32chars!",
			},
			auth,
			withHeader("X-Admin-Id", testAdminID),
		)
		drain(resp)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("GetClient", func(t *testing.T) {
		resp := do(http.MethodGet,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/clients/%s", userID, clientID),
			nil, auth,
		)
		var body map[string]any
		decodeBody(resp, &body)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, clientID, body["id"])
	})

	t.Run("AddClientAllowedFlow", func(t *testing.T) {
		resp := do(http.MethodPost,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/clients/%s/flows", userID, clientID),
			map[string]any{"flow": "authorization_code"},
			auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("AddRedirectUri", func(t *testing.T) {
		resp := do(http.MethodPost,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/clients/%s/redirect-uris", userID, clientID),
			map[string]any{"redirect_uri": "https://example.com/callback"},
			auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("AddScopeAllowedClient", func(t *testing.T) {
		resp := do(http.MethodPost,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/protected-resources/%s/scopes/%s/allowed-clients/%s",
				userID, prID, scopeID, clientID),
			map[string]any{"client_id": clientID},
			auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("RemoveScopeAllowedClient", func(t *testing.T) {
		resp := do(http.MethodDelete,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/protected-resources/%s/scopes/%s/allowed-clients/%s",
				userID, prID, scopeID, clientID),
			nil, auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("SoftDeleteScope", func(t *testing.T) {
		resp := do(http.MethodDelete,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/protected-resources/%s/scopes/%s", userID, prID, scopeID),
			nil, auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("SoftDeleteClient", func(t *testing.T) {
		resp := do(http.MethodDelete,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/clients/%s", userID, clientID),
			nil, auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("SoftDeleteProtectedResource", func(t *testing.T) {
		resp := do(http.MethodDelete,
			fmt.Sprintf("/v1/backoffice/admins/users/%s/protected-resources/%s", userID, prID),
			nil, auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("SoftDeleteUser", func(t *testing.T) {
		resp := do(http.MethodDelete,
			fmt.Sprintf("/v1/backoffice/admins/users/%s", userID),
			nil, auth,
		)
		drain(resp)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
