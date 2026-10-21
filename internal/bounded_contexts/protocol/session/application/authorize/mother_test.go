package authorize_test

import (
	"context"

	"github.com/google/uuid"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
)

type mockClientRepo struct {
	client *session.Client
	err    error
}

func (m *mockClientRepo) FindById(_ context.Context, _ string) (*session.Client, error) {
	return m.client, m.err
}

type mockIdpFinderRepo struct {
	idp *protocolShared.IdentityProviderConfig
	err error
}

func (m *mockIdpFinderRepo) FindByReferenceName(_ context.Context, _ string) (*protocolShared.IdentityProviderConfig, error) {
	return m.idp, m.err
}

type mockIdpClient struct {
	url string
	err error
}

func (m *mockIdpClient) BuildAuthorizationURL(_ *protocolShared.IdentityProviderConfig, _ string) (string, error) {
	return m.url, m.err
}

type mockStateCodec struct {
	encoded string
	err     error
}

func (m *mockStateCodec) Encode(_ session.IdpStatePayload) (string, error) {
	return m.encoded, m.err
}

func (m *mockStateCodec) Decode(_ string) (session.IdpStatePayload, error) {
	return session.IdpStatePayload{}, nil
}

func newClientMother(clientID string) *session.Client {
	return session.NewClient(
		clientID, "Test Client",
		uuid.New().String(),
		[]string{"authorization_code"},
		[]string{"https://example.com/callback"},
		[]*session.Scope{{ScopeName: "openid"}},
		false,
	)
}

func aGoogleIdp() *protocolShared.IdentityProviderConfig {
	return protocolShared.NewIdentityProviderConfig(
		"idp-id", "Google", "https://accounts.google.com", "client-id", "",
		"https://accounts.google.com/o/oauth2/v2/auth", "", "https://www.googleapis.com/oauth2/v3/certs",
		[]string{"openid"},
		"https://api.example.com/v1/auth/callback/google",
	)
}
