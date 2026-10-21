package authenticate_client_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
)

func publicClientMother(clientID string, public bool) *session.Client {
	return session.NewClient(clientID, "n", uuid.New().String(), nil, nil, nil, public)
}

func newPublicClientMethod(client *session.Client, err error) *authenticate_client.PublicClientMethod {
	return authenticate_client.NewPublicClientMethod(
		domain_services.NewClientFinder(&mockPublicClientRepo{client: client, err: err}),
	)
}

func TestPublicClientMethod_CanProcess(t *testing.T) {
	t.Parallel()

	m := newPublicClientMethod(nil, nil)

	tests := []struct {
		name  string
		creds authenticate_client.Credentials
		want  bool
	}{
		{name: "client_id_only", creds: authenticate_client.Credentials{IsFormEncoded: true, ClientID: "cid"}, want: true},
		{name: "client_id_with_secret", creds: authenticate_client.Credentials{IsFormEncoded: true, ClientID: "cid", ClientSecret: "s"}, want: false},
		{name: "no_client_id", creds: authenticate_client.Credentials{IsFormEncoded: true}, want: false},
		{name: "not_form_encoded", creds: authenticate_client.Credentials{ClientID: "cid"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, m.CanProcess(tt.creds))
		})
	}
}

func TestPublicClientMethod_Authenticate(t *testing.T) {
	t.Parallel()

	t.Run("public client returns client id", func(t *testing.T) {
		t.Parallel()
		m := newPublicClientMethod(publicClientMother("pub-cid", true), nil)
		result, err := m.Authenticate(context.Background(), authenticate_client.Credentials{ClientID: "pub-cid"})
		require.NoError(t, err)
		assert.Equal(t, "pub-cid", result.ClientID)
		assert.True(t, result.IsPublic)
	})

	t.Run("confidential client returns error", func(t *testing.T) {
		t.Parallel()
		m := newPublicClientMethod(publicClientMother("conf-cid", false), nil)
		_, err := m.Authenticate(context.Background(), authenticate_client.Credentials{ClientID: "conf-cid"})
		assert.Error(t, err)
	})

	t.Run("finder error propagates", func(t *testing.T) {
		t.Parallel()
		m := newPublicClientMethod(nil, assert.AnError)
		_, err := m.Authenticate(context.Background(), authenticate_client.Credentials{ClientID: "cid"})
		assert.Error(t, err)
	})
}
