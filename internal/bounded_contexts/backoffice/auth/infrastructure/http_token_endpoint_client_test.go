package infrastructure_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/request_id"
)

func TestHttpTokenEndpointClient_ExchangeCode(t *testing.T) {
	t.Parallel()

	t.Run("forwards the request id from the context as a header", func(t *testing.T) {
		t.Parallel()
		var capturedHeader string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedHeader = r.Header.Get(request_id.Header)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "signed.jwt",
				"id_token":     "signed.id.jwt",
				"token_type":   "Bearer",
				"expires_in":   3600,
			})
		}))
		defer server.Close()

		client := infrastructure.NewHttpTokenEndpointClient(server.URL, "client-id", "client-secret", "https://example.com/cb")
		ctx := request_id.WithContext(context.Background(), "the-request-id")

		result, err := client.ExchangeCode(ctx, "code", "verifier")

		require.NoError(t, err)
		assert.Equal(t, "signed.jwt", result.AccessToken)
		assert.Equal(t, "signed.id.jwt", result.IDToken)
		assert.Equal(t, "Bearer", result.TokenType)
		assert.Equal(t, 3600, result.ExpiresIn)
		assert.Equal(t, "the-request-id", capturedHeader)
	})

	t.Run("omits the header when the context carries no request id", func(t *testing.T) {
		t.Parallel()
		var headerWasSet bool
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, headerWasSet = r.Header[request_id.Header]
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"access_token": "signed.jwt"})
		}))
		defer server.Close()

		client := infrastructure.NewHttpTokenEndpointClient(server.URL, "client-id", "client-secret", "https://example.com/cb")

		_, err := client.ExchangeCode(context.Background(), "code", "verifier")

		require.NoError(t, err)
		assert.False(t, headerWasSet)
	})
}
