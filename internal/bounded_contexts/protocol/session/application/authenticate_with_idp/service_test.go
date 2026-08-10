package authenticate_with_idp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	authn "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_with_idp"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func newService(sessionRepo *mockSessionRepo, idpRepo *mockIdpRepo, client *mockIdpClient) *authn.Service {
	factory := domain_services.NewIdentityProviderClientFactory(map[string]session.IdentityProviderClient{
		"google": client,
	})
	return authn.NewService(
		domain_services.NewSessionFinder(sessionRepo),
		domain_services.NewIdentityProviderFinder(idpRepo),
		factory,
	)
}

func TestAuthenticate_HappyPath(t *testing.T) {
	t.Parallel()

	svc := newService(
		&mockSessionRepo{session: aSession(t)},
		&mockIdpRepo{idp: aGoogleIDP()},
		&mockIdpClient{url: "https://accounts.google.com/o/oauth2/v2/auth?state=x"},
	)

	result, err := svc.Authenticate(context.Background(), authn.NewQuery(ids.NewSessionId(), "idp-123"))
	require.NoError(t, err)
	assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth?state=x", result)
}

func TestAuthenticate_SessionNotFound(t *testing.T) {
	t.Parallel()

	svc := newService(
		&mockSessionRepo{session: nil},
		&mockIdpRepo{idp: aGoogleIDP()},
		&mockIdpClient{url: "x"},
	)

	_, err := svc.Authenticate(context.Background(), authn.NewQuery(ids.NewSessionId(), "idp-123"))
	assert.Error(t, err)
}

func TestAuthenticate_IdpNotFound(t *testing.T) {
	t.Parallel()

	svc := newService(
		&mockSessionRepo{session: aSession(t)},
		&mockIdpRepo{idp: nil},
		&mockIdpClient{url: "x"},
	)

	_, err := svc.Authenticate(context.Background(), authn.NewQuery(ids.NewSessionId(), "missing"))
	assert.Error(t, err)
}

func TestAuthenticate_UnsupportedIdp(t *testing.T) {
	t.Parallel()

	unsupported := session.NewIdentityProvider(
		"idp-999", "github", "https://github.com", "cid", "secret",
		"https://github.com/login/oauth/authorize", "https://github.com/login/oauth/access_token",
		"", "", []string{"read:user"},
	)

	svc := newService(
		&mockSessionRepo{session: aSession(t)},
		&mockIdpRepo{idp: unsupported},
		&mockIdpClient{url: "x"},
	)

	_, err := svc.Authenticate(context.Background(), authn.NewQuery(ids.NewSessionId(), "idp-999"))
	assert.Error(t, err)
}

func TestAuthenticate_IdpClientError(t *testing.T) {
	t.Parallel()

	svc := newService(
		&mockSessionRepo{session: aSession(t)},
		&mockIdpRepo{idp: aGoogleIDP()},
		&mockIdpClient{err: errBoom},
	)

	_, err := svc.Authenticate(context.Background(), authn.NewQuery(ids.NewSessionId(), "idp-123"))
	assert.Error(t, err)
}
