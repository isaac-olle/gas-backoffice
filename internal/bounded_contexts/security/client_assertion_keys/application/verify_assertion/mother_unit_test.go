package verify_assertion_test

import (
	"context"
)

// mockAssertionVerifier stubs client_assertion_keys.AssertionVerifier.
type mockAssertionVerifier struct {
	claims map[string]any
	err    error
}

func (m *mockAssertionVerifier) Verify(_ context.Context, _ string) (map[string]any, error) {
	return m.claims, m.err
}
