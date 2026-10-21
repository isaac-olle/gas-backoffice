package authenticate_client

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

const jwtBearerAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"

var ErrInvalidAssertion app_errors.ErrorCode = "INVALID_CLIENT_ASSERTION"

// ClientAssertionVerifier verifies a client_assertion JWT and returns the
// client id it asserts. It is satisfied by an adapter over the security
// bounded context, since assertion verification is cross-BC.
type ClientAssertionVerifier interface {
	Verify(ctx context.Context, assertion string) (clientID string, err error)
}

// JWTAssertionMethod implements private_key_jwt (RFC 7523).
type JWTAssertionMethod struct {
	verifier ClientAssertionVerifier
}

func NewJWTAssertionMethod(verifier ClientAssertionVerifier) *JWTAssertionMethod {
	return &JWTAssertionMethod{verifier: verifier}
}

func (m *JWTAssertionMethod) CanProcess(creds Credentials) bool {
	return creds.ClientAssertionType == jwtBearerAssertionType
}

func (m *JWTAssertionMethod) Authenticate(ctx context.Context, creds Credentials) (Result, error) {
	if creds.ClientAssertion == "" {
		return Result{}, app_errors.UnauthorizedError(ErrInvalidAssertion).
			WithDescription("client_assertion is required when client_assertion_type is jwt-bearer")
	}

	clientID, err := m.verifier.Verify(ctx, creds.ClientAssertion)
	if err != nil {
		return Result{}, err
	}
	if clientID == "" {
		return Result{}, app_errors.UnauthorizedError(ErrInvalidAssertion).
			WithDescription("invalid client assertion claims")
	}
	return Result{ClientID: clientID, IsPublic: false}, nil
}
