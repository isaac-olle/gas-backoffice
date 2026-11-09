package infrastructure_test

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure/helpers/oidc"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/stretchr/testify/require"
)

type stateTestEncrypter struct{ err error }

func (e stateTestEncrypter) Encrypt(s string) (string, error) {
	if e.err != nil {
		return "", e.err
	}
	return base64.StdEncoding.EncodeToString([]byte(s)), nil
}
func (e stateTestEncrypter) Decrypt(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	return string(b), err
}
func (e stateTestEncrypter) Seal(b []byte) ([]byte, error) {
	if e.err != nil {
		return nil, e.err
	}
	return b, nil
}
func (e stateTestEncrypter) Open(b []byte) ([]byte, error) { return b, nil }

var _ persistance.Encrypter = stateTestEncrypter{}

func newAuthorizer(endpoint string, encryptErr error) *infrastructure.OidcAuthorizer {
	encoder := oidc.NewStateEncoder(5*time.Minute, stateTestEncrypter{err: encryptErr})
	builder := oidc.NewRequestBuilder(oidc.NewRandomCodeVerifierGenerator(), encoder)
	return infrastructure.NewOidcAuthorizer(builder, endpoint, "backoffice-client", "https://app.example/callback")
}

func TestOidcAuthorizerConnect(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		location    string
		encryptErr  error
		unreachable bool
		wantErr     bool
	}{
		{name: "returns the upstream redirect location", status: http.StatusFound, location: "https://idp.example/authorize?state=xyz"},
		{name: "rejects upstream error responses", status: http.StatusBadRequest, location: "", wantErr: true},
		{name: "rejects redirects without a Location", status: http.StatusFound, wantErr: true},
		{name: "wraps state encryption failures", status: http.StatusFound, encryptErr: errors.New("encryption failed"), wantErr: true},
		{name: "reports an unreachable authorization endpoint", unreachable: true, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := "http://127.0.0.1:1/authorize"
			var received *http.Request
			if !tt.unreachable {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					received = r
					if tt.location != "" {
						w.Header().Set("Location", tt.location)
					}
					w.WriteHeader(tt.status)
				}))
				t.Cleanup(server.Close)
				endpoint = server.URL
			}
			got, err := newAuthorizer(endpoint, tt.encryptErr).Connect(context.Background(), &auth.IdpConfig{IdpReferenceName: "google"})
			require.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.location, got)
			require.NotNil(t, received)
			require.Equal(t, http.MethodGet, received.Method)
			query := received.URL.Query()
			require.Equal(t, "code", query.Get("response_type"))
			require.Equal(t, "backoffice-client", query.Get("client_id"))
			require.Equal(t, "https://app.example/callback", query.Get("redirect_uri"))
			require.Equal(t, "google", query.Get("idp_reference_name"))
			require.Equal(t, "S256", query.Get("code_challenge_method"))
			require.NotEmpty(t, query.Get("code_challenge"))
			require.NotEmpty(t, query.Get("state"))
			require.False(t, strings.Contains(query.Get("state"), "+"))
		})
	}
}
