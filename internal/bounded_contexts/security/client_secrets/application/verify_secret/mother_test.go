package verify_secret_test

import (
	"context"
)

// mockSecretVerifier stubs client_secrets.SecretVerifier.
type mockSecretVerifier struct {
	err error
}

func (m *mockSecretVerifier) Verify(_ context.Context, _, _ string) error { return m.err }
