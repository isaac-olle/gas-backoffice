package infrastructure_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure/helpers/oidc"
	"github.com/isaac-olle/gas-backoffice/internal/shared/request_id"
	"github.com/stretchr/testify/require"
)

func TestHttpTokenEndpointClientExchangeCode(t *testing.T) {
	tests := []struct {
		name          string
		status        int
		body          string
		requestID     string
		wantRequestID string
		wantError     string
	}{
		{name: "forwards request id and authorization code form", status: http.StatusOK, body: `{"id_token":"invalid.jwt"}`, requestID: "request-123", wantRequestID: "request-123", wantError: "verify id_token"},
		{name: "omits request id when absent", status: http.StatusOK, body: `{"id_token":"invalid.jwt"}`, wantError: "verify id_token"},
		{name: "reports non-success response", status: http.StatusBadRequest, body: `{"error":"invalid_grant"}`, wantError: "invalid_grant"},
		{name: "reports malformed response json", status: http.StatusOK, body: `{`, wantError: "decode token response"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var request *http.Request
			var form url.Values
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				request = r
				if err := r.ParseForm(); err == nil {
					form = r.PostForm
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()
			verifier := infrastructure.NewGasIdTokenVerifier("http://127.0.0.1:1/jwks", "issuer", "client-id")
			client := infrastructure.NewHttpTokenEndpointClient(server.URL, "client-id", "client-secret", "https://app.example/callback", verifier)
			ctx := context.Background()
			if tt.requestID != "" {
				ctx = request_id.WithContext(ctx, tt.requestID)
			}
			_, err := client.ExchangeCode(ctx, "auth-code", "code-verifier", "google")
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantError)
			require.NotNil(t, request)
			require.Equal(t, http.MethodPost, request.Method)
			user, password, ok := request.BasicAuth()
			require.True(t, ok)
			require.Equal(t, "client-id", user)
			require.Equal(t, "client-secret", password)
			if tt.wantRequestID != "" {
				require.Equal(t, tt.wantRequestID, request.Header.Get(request_id.Header))
			} else {
				require.Empty(t, request.Header.Get(request_id.Header))
			}
			if form != nil && tt.status == http.StatusOK && strings.Contains(tt.body, "id_token") {
				require.Equal(t, "authorization_code", form.Get("grant_type"))
				require.Equal(t, "auth-code", form.Get("code"))
				require.Equal(t, "code-verifier", form.Get("code_verifier"))
				require.Equal(t, "https://app.example/callback", form.Get("redirect_uri"))
			}
		})
	}
}

func TestTokenRedirectMatchesAuthorizationRedirect(t *testing.T) {
	encoder := oidc.NewStateEncoder(time.Minute, stateTestEncrypter{})
	builder := oidc.NewRequestBuilder(oidc.NewRandomCodeVerifierGenerator(), encoder)
	const redirectTemplate = "https://app.example/callback/:idp_reference_name"
	authorizationURL, err := builder.Build(redirectTemplate, "client", "google", "https://gas.example/oidc/authorize")
	require.NoError(t, err)
	parsed, err := url.Parse(authorizationURL)
	require.NoError(t, err)
	state, err := encoder.DecodeState(parsed.Query().Get("state"))
	require.NoError(t, err)
	require.Equal(t, "google", state.IdpReferenceName)
	var receivedRedirect string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		receivedRedirect = r.PostForm.Get("redirect_uri")
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()
	client := infrastructure.NewHttpTokenEndpointClient(server.URL, "client", "secret", redirectTemplate, nil)
	_, err = client.ExchangeCode(context.Background(), "code", state.CodeVerifier, state.IdpReferenceName)
	require.Error(t, err)
	require.Equal(t, "https://app.example/callback/google", receivedRedirect)
	require.Equal(t, parsed.Query().Get("redirect_uri"), receivedRedirect)
}
