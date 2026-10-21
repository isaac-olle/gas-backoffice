package authenticate_client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
)

func TestBodyCredentialsMethod_CanProcess(t *testing.T) {
	t.Parallel()

	m := authenticate_client.NewBodyCredentialsMethod(&stubSecretVerifier{})

	tests := []struct {
		name  string
		creds authenticate_client.Credentials
		want  bool
	}{
		{name: "not_form_encoded", creds: authenticate_client.Credentials{ClientID: "cid", ClientSecret: "s"}, want: false},
		{name: "no_secret", creds: authenticate_client.Credentials{IsFormEncoded: true, ClientID: "cid"}, want: false},
		{name: "no_client_id", creds: authenticate_client.Credentials{IsFormEncoded: true, ClientSecret: "s"}, want: false},
		{name: "id_and_secret", creds: authenticate_client.Credentials{IsFormEncoded: true, ClientID: "cid", ClientSecret: "s"}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, m.CanProcess(tt.creds))
		})
	}
}

func TestBodyCredentialsMethod_Authenticate(t *testing.T) {
	t.Parallel()

	t.Run("verifier error propagates", func(t *testing.T) {
		t.Parallel()
		m := authenticate_client.NewBodyCredentialsMethod(&stubSecretVerifier{err: assert.AnError})
		_, err := m.Authenticate(context.Background(), authenticate_client.Credentials{ClientID: "cid", ClientSecret: "s"})
		assert.Error(t, err)
	})

	t.Run("success returns client id, not public", func(t *testing.T) {
		t.Parallel()
		m := authenticate_client.NewBodyCredentialsMethod(&stubSecretVerifier{})
		result, err := m.Authenticate(context.Background(), authenticate_client.Credentials{ClientID: "cid", ClientSecret: "s"})
		require.NoError(t, err)
		assert.Equal(t, "cid", result.ClientID)
		assert.False(t, result.IsPublic)
	})
}
