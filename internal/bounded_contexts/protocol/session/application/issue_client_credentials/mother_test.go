package issue_client_credentials_test

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/buses/event"

	"github.com/google/uuid"

	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session"
)

type mockClientRepo struct {
	client *session.Client
	err    error
}

func (m *mockClientRepo) FindById(_ context.Context, _ string) (*session.Client, error) {
	return m.client, m.err
}

type mockTokenGen struct {
	token string
	err   error
}

func (m *mockTokenGen) GenerateFromSession(_ context.Context, _ *session.AuthorizationCodeSession, _ *session.Client) (string, error) {
	return m.token, m.err
}
func (m *mockTokenGen) GenerateFromClientCredentials(_ context.Context, _ *session.ClientCredentialsSession) (string, error) {
	return m.token, m.err
}

type mockEventBus struct {
	published     [][]event.Event
	publishCalled bool
}

func (m *mockEventBus) Publish(_ context.Context, evts []event.Event) {
	m.publishCalled = true
	m.published = append(m.published, evts)
}

func newClientWithCCMother(clientID string) *session.Client {
	return session.NewClient(
		clientID, "TestClient",
		uuid.New().String(), uuid.New().String(),
		[]string{"client_credentials"},
		[]string{},
		[]*session.Scope{{ScopeName: "openid"}},
	)
}
