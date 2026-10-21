package authenticate_client

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var ErrNoAuthenticationMethod app_errors.ErrorCode = "NO_AUTHENTICATION_METHOD"

// Credentials is a neutral view of whatever authentication material the
// /oidc/token request carried (Authorization header, form fields). Building
// it is the HTTP layer's job; deciding what it means is this package's job.
type Credentials struct {
	AuthorizationHeader string
	IsFormEncoded       bool
	ClientID            string
	ClientSecret        string
	ClientAssertionType string
	ClientAssertion     string
}

// Result is the outcome of a successful authentication: who the client is
// and whether it authenticated as a public client.
type Result struct {
	ClientID string
	IsPublic bool
}

// Method resolves one specific way a client can authenticate at /oidc/token.
type Method interface {
	CanProcess(creds Credentials) bool
	Authenticate(ctx context.Context, creds Credentials) (Result, error)
}

// Service picks the first Method that matches the given credentials and
// delegates to it. Order matters: methods are tried in the order given to
// NewService.
type Service struct {
	methods []Method
}

func NewService(methods ...Method) *Service {
	return &Service{methods: methods}
}

func (s *Service) Authenticate(ctx context.Context, creds Credentials) (Result, error) {
	for _, m := range s.methods {
		if !m.CanProcess(creds) {
			continue
		}
		return m.Authenticate(ctx, creds)
	}
	return Result{}, app_errors.UnauthorizedError(ErrNoAuthenticationMethod).
		WithDescription("no supported authentication method found in request")
}
