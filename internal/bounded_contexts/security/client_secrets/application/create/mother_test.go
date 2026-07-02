package create_test

// mockSecretHasher stubs client_secrets.SecretHasher.
type mockSecretHasher struct {
	hash string
	err  error
}

func (m *mockSecretHasher) Hash(_ string) (string, error) { return m.hash, m.err }
