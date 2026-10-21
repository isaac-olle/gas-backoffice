package authenticate_client

import "context"

// BodyCredentialsMethod implements client_secret_post (RFC 6749 §2.3.1).
type BodyCredentialsMethod struct {
	verifier ClientSecretVerifier
}

func NewBodyCredentialsMethod(verifier ClientSecretVerifier) *BodyCredentialsMethod {
	return &BodyCredentialsMethod{verifier: verifier}
}

func (m *BodyCredentialsMethod) CanProcess(creds Credentials) bool {
	return creds.IsFormEncoded && creds.ClientID != "" && creds.ClientSecret != ""
}

func (m *BodyCredentialsMethod) Authenticate(ctx context.Context, creds Credentials) (Result, error) {
	if err := m.verifier.Verify(ctx, creds.ClientID, creds.ClientSecret); err != nil {
		return Result{}, err
	}
	return Result{ClientID: creds.ClientID, IsPublic: false}, nil
}
