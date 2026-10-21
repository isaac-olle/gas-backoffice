package authenticate_client_test

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
)

type stubSecretVerifier struct {
	err              error
	receivedClientID string
	receivedSecret   string
}

func (v *stubSecretVerifier) Verify(_ context.Context, clientID, secret string) error {
	v.receivedClientID = clientID
	v.receivedSecret = secret
	return v.err
}

type stubAssertionVerifier struct {
	clientID string
	err      error
}

func (v *stubAssertionVerifier) Verify(_ context.Context, _ string) (string, error) {
	return v.clientID, v.err
}

type mockPublicClientRepo struct {
	client *session.Client
	err    error
}

func (m *mockPublicClientRepo) FindById(_ context.Context, _ string) (*session.Client, error) {
	return m.client, m.err
}

type stubMethod struct {
	canProcess bool
	result     authenticate_client.Result
	err        error
}

func (m *stubMethod) CanProcess(_ authenticate_client.Credentials) bool { return m.canProcess }
func (m *stubMethod) Authenticate(_ context.Context, _ authenticate_client.Credentials) (authenticate_client.Result, error) {
	return m.result, m.err
}
