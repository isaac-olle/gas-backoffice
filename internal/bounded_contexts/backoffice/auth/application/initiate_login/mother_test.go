package initiate_login_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
)

type stubCodeVerifierGenerator struct {
	value string
	err   error
}

func (g stubCodeVerifierGenerator) GenerateCodeVerifier() (vos.CodeVerifier, error) {
	if g.err != nil {
		return vos.CodeVerifier{}, g.err
	}
	return vos.NewCodeVerifierFromString(g.value)
}

func aValidVerifier() string { return strings.Repeat("a", 43) }

type stubIdpConfigRepo struct {
	idp *auth.IdpConfig
	err error

	gotReferenceName vos.IdpReferenceName
}

func (r *stubIdpConfigRepo) FindByReferenceName(
	_ context.Context,
	referenceName vos.IdpReferenceName,
) (*auth.IdpConfig, error) {
	r.gotReferenceName = referenceName
	return r.idp, r.err
}

type stubAuthorizer struct {
	url string
	err error

	gotIdp          auth.IdpConfig
	gotLoginRequest *auth.LoginRequest
}

func (a *stubAuthorizer) Authorize(
	_ context.Context,
	idp auth.IdpConfig,
	loginRequest *auth.LoginRequest,
) (string, error) {
	a.gotIdp = idp
	a.gotLoginRequest = loginRequest
	return a.url, a.err
}

func anIdp() *auth.IdpConfig {
	return &auth.IdpConfig{}
}

var errBoom = fmt.Errorf("boom")
