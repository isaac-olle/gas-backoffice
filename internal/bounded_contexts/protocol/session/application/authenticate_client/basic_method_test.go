package authenticate_client_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
)

func basicHeader(clientID, secret string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+secret))
}

func TestBasicMethod_CanProcess(t *testing.T) {
	t.Parallel()

	m := authenticate_client.NewBasicMethod(&stubSecretVerifier{})

	tests := []struct {
		name  string
		creds authenticate_client.Credentials
		want  bool
	}{
		{name: "no_authorization_header", creds: authenticate_client.Credentials{}, want: false},
		{name: "bearer_token_header", creds: authenticate_client.Credentials{AuthorizationHeader: "Bearer sometoken"}, want: false},
		{name: "basic_header", creds: authenticate_client.Credentials{AuthorizationHeader: basicHeader("abc", "secret")}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, m.CanProcess(tt.creds))
		})
	}
}

func TestBasicMethod_Authenticate(t *testing.T) {
	t.Parallel()

	t.Run("invalid base64 returns error", func(t *testing.T) {
		t.Parallel()
		m := authenticate_client.NewBasicMethod(&stubSecretVerifier{})
		_, err := m.Authenticate(context.Background(), authenticate_client.Credentials{AuthorizationHeader: "Basic not-valid-base64!!!"})
		assert.Error(t, err)
	})

	t.Run("no colon in credentials returns error", func(t *testing.T) {
		t.Parallel()
		m := authenticate_client.NewBasicMethod(&stubSecretVerifier{})
		header := "Basic " + base64.StdEncoding.EncodeToString([]byte("justid"))
		_, err := m.Authenticate(context.Background(), authenticate_client.Credentials{AuthorizationHeader: header})
		assert.Error(t, err)
	})

	t.Run("verifier error propagates", func(t *testing.T) {
		t.Parallel()
		m := authenticate_client.NewBasicMethod(&stubSecretVerifier{err: assert.AnError})
		_, err := m.Authenticate(context.Background(), authenticate_client.Credentials{AuthorizationHeader: basicHeader("abc", "secret")})
		assert.Error(t, err)
	})

	t.Run("success returns client id, not public", func(t *testing.T) {
		t.Parallel()
		verifier := &stubSecretVerifier{}
		m := authenticate_client.NewBasicMethod(verifier)
		result, err := m.Authenticate(context.Background(), authenticate_client.Credentials{AuthorizationHeader: basicHeader("abc", "secret")})
		require.NoError(t, err)
		assert.Equal(t, "abc", result.ClientID)
		assert.False(t, result.IsPublic)
		assert.Equal(t, "abc", verifier.receivedClientID)
		assert.Equal(t, "secret", verifier.receivedSecret)
	})
}
