package authenticate_client

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var ErrInvalidBasicAuth app_errors.ErrorCode = "INVALID_BASIC_AUTH"

// ClientSecretVerifier verifies a client_id/client_secret pair. It is
// satisfied by an adapter over the security bounded context, since secret
// verification is cross-BC.
type ClientSecretVerifier interface {
	Verify(ctx context.Context, clientID, secret string) error
}

// BasicMethod implements client_secret_basic (RFC 6749 §2.3.1).
type BasicMethod struct {
	verifier ClientSecretVerifier
}

func NewBasicMethod(verifier ClientSecretVerifier) *BasicMethod {
	return &BasicMethod{verifier: verifier}
}

func (m *BasicMethod) CanProcess(creds Credentials) bool {
	return strings.HasPrefix(creds.AuthorizationHeader, "Basic ")
}

func (m *BasicMethod) Authenticate(ctx context.Context, creds Credentials) (Result, error) {
	encoded := strings.TrimPrefix(creds.AuthorizationHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return Result{}, app_errors.UnauthorizedError(ErrInvalidBasicAuth).
			WithDescription("invalid Basic Auth encoding")
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return Result{}, app_errors.UnauthorizedError(ErrInvalidBasicAuth).
			WithDescription("invalid Basic Auth format: expected client_id:client_secret")
	}

	clientID, secret := parts[0], parts[1]
	if err = m.verifier.Verify(ctx, clientID, secret); err != nil {
		return Result{}, err
	}
	return Result{ClientID: clientID, IsPublic: false}, nil
}
