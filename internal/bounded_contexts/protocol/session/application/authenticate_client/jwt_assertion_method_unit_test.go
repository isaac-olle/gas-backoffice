package authenticate_client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
)

func TestJWTAssertionMethod_CanProcess(t *testing.T) {
	t.Parallel()

	m := authenticate_client.NewJWTAssertionMethod(&stubAssertionVerifier{})

	tests := []struct {
		name  string
		creds authenticate_client.Credentials
		want  bool
	}{
		{name: "no_assertion_type", creds: authenticate_client.Credentials{}, want: false},
		{name: "wrong_assertion_type", creds: authenticate_client.Credentials{ClientAssertionType: "something-else"}, want: false},
		{name: "jwt_bearer_assertion_type", creds: authenticate_client.Credentials{ClientAssertionType: "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, m.CanProcess(tt.creds))
		})
	}
}

func TestJWTAssertionMethod_Authenticate(t *testing.T) {
	t.Parallel()

	t.Run("missing assertion returns error", func(t *testing.T) {
		t.Parallel()
		m := authenticate_client.NewJWTAssertionMethod(&stubAssertionVerifier{})
		_, err := m.Authenticate(context.Background(), authenticate_client.Credentials{})
		assert.Error(t, err)
	})

	t.Run("verifier error propagates", func(t *testing.T) {
		t.Parallel()
		m := authenticate_client.NewJWTAssertionMethod(&stubAssertionVerifier{err: assert.AnError})
		_, err := m.Authenticate(context.Background(), authenticate_client.Credentials{ClientAssertion: "jwt"})
		assert.Error(t, err)
	})

	t.Run("empty client id from verifier returns error", func(t *testing.T) {
		t.Parallel()
		m := authenticate_client.NewJWTAssertionMethod(&stubAssertionVerifier{clientID: ""})
		_, err := m.Authenticate(context.Background(), authenticate_client.Credentials{ClientAssertion: "jwt"})
		assert.Error(t, err)
	})

	t.Run("success returns client id, not public", func(t *testing.T) {
		t.Parallel()
		m := authenticate_client.NewJWTAssertionMethod(&stubAssertionVerifier{clientID: "abc"})
		result, err := m.Authenticate(context.Background(), authenticate_client.Credentials{ClientAssertion: "jwt"})
		require.NoError(t, err)
		assert.Equal(t, "abc", result.ClientID)
		assert.False(t, result.IsPublic)
	})
}
