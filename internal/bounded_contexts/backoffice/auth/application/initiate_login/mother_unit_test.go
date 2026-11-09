package initiate_login_test

import (
	"context"
	"errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
)

type idpConfigRepositoryStub struct {
	idp     *auth.IdpConfig
	err     error
	gotName vos.IdpReferenceName
}

func (r *idpConfigRepositoryStub) FindByReferenceName(_ context.Context, name vos.IdpReferenceName) (*auth.IdpConfig, error) {
	r.gotName = name
	return r.idp, r.err
}

type idpConnectorStub struct {
	url string
	err error
	got *auth.IdpConfig
}

func (c *idpConnectorStub) Connect(_ context.Context, idp *auth.IdpConfig) (string, error) {
	c.got = idp
	return c.url, c.err
}

func validIdp() *auth.IdpConfig { return &auth.IdpConfig{IdpReferenceName: "google"} }

var errInitiateLogin = errors.New("dependency failure")
