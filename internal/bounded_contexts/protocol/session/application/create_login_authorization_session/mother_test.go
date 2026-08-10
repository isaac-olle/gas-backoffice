package create_login_authorization_session_test

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type mockSessionRepo struct {
	saved []*session.AuthorizationCodeSession
	err   error
}

func (m *mockSessionRepo) Save(_ context.Context, s *session.AuthorizationCodeSession) error {
	if m.err != nil {
		return m.err
	}
	m.saved = append(m.saved, s)
	return nil
}
func (m *mockSessionRepo) FindById(_ context.Context, _ ids.SessionId) (*session.AuthorizationCodeSession, error) {
	return nil, nil
}
func (m *mockSessionRepo) FindByCode(_ context.Context, _ string) (*session.AuthorizationCodeSession, error) {
	return nil, nil
}
func (m *mockSessionRepo) SoftDelete(_ context.Context, _ ids.SessionId) error { return nil }
func (m *mockSessionRepo) HardDelete(_ context.Context, _ ids.SessionId) error { return nil }

type mockEventBus struct {
	publishCalled bool
}

func (m *mockEventBus) Publish(_ context.Context, _ []event.Event) { m.publishCalled = true }

type mockCodeGen struct{ code string }

func (c *mockCodeGen) Generate() (string, error) { return c.code, nil }
