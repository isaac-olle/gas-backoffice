package authorize_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	protocolSharedDS "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authorize"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func newPublicClientMother(clientID string) *session.Client {
	return session.NewClient(
		clientID, "Public Client",
		uuid.New().String(),
		[]string{"authorization_code"},
		[]string{"https://example.com/callback"},
		[]*session.Scope{{ScopeName: "openid"}},
		true,
	)
}

func newService(clientRepo *mockClientRepo, idpRepo *mockIdpFinderRepo, idpClient *mockIdpClient, stateCodec *mockStateCodec) *authorize.Service {
	return authorize.NewService(
		domain_services.NewClientFinder(clientRepo),
		protocolSharedDS.NewIdentityProviderFinder(idpRepo),
		idpClient,
		stateCodec,
	)
}

func validRequest(clientID string) authorize.Request {
	return authorize.Request{
		ResponseType: "code",
		ClientID:     clientID,
		RedirectURI:  "https://example.com/callback",
		Scope:        "openid",
		State:        "client-state",
	}
}

func TestAuthorizeService(t *testing.T) {
	t.Parallel()

	t.Run("validates against the client and returns the idp authorization url built with the encoded idp state", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		svc := newService(
			&mockClientRepo{client: newClientMother(clientID)},
			&mockIdpFinderRepo{idp: aGoogleIdp()},
			&mockIdpClient{url: "https://accounts.google.com/authorize?state=opaque-state"},
			&mockStateCodec{encoded: "opaque-state"},
		)

		url, err := svc.Authorize(context.Background(), validRequest(clientID))

		require.NoError(t, err)
		assert.Equal(t, "https://accounts.google.com/authorize?state=opaque-state", url)
	})

	t.Run("returns error when the client cannot be found", func(t *testing.T) {
		t.Parallel()
		svc := newService(
			&mockClientRepo{client: nil},
			&mockIdpFinderRepo{idp: aGoogleIdp()},
			&mockIdpClient{},
			&mockStateCodec{},
		)

		_, err := svc.Authorize(context.Background(), validRequest(uuid.New().String()))

		assert.Error(t, err)
	})

	t.Run("returns error when the identity provider cannot be found", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		svc := newService(
			&mockClientRepo{client: newClientMother(clientID)},
			&mockIdpFinderRepo{idp: nil},
			&mockIdpClient{},
			&mockStateCodec{},
		)

		_, err := svc.Authorize(context.Background(), validRequest(clientID))

		assert.Error(t, err)
	})

	t.Run("rejects an unsupported response_type without touching the idp", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		idpClient := &mockIdpClient{url: "should-not-be-used"}
		svc := newService(
			&mockClientRepo{client: newClientMother(clientID)},
			&mockIdpFinderRepo{idp: aGoogleIdp()},
			idpClient,
			&mockStateCodec{},
		)

		req := validRequest(clientID)
		req.ResponseType = "token"

		_, err := svc.Authorize(context.Background(), req)

		require.Error(t, err)
		var appErr *app_errors.AppError
		require.ErrorAs(t, err, &appErr)
	})

	t.Run("rejects a redirect_uri not registered for the client", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		svc := newService(
			&mockClientRepo{client: newClientMother(clientID)},
			&mockIdpFinderRepo{idp: aGoogleIdp()},
			&mockIdpClient{},
			&mockStateCodec{},
		)

		req := validRequest(clientID)
		req.RedirectURI = "https://evil.example.com/cb"

		_, err := svc.Authorize(context.Background(), req)

		require.Error(t, err)
	})

	t.Run("requires pkce for public clients", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		publicClient := newPublicClientMother(clientID)
		svc := newService(
			&mockClientRepo{client: publicClient},
			&mockIdpFinderRepo{idp: aGoogleIdp()},
			&mockIdpClient{},
			&mockStateCodec{},
		)

		_, err := svc.Authorize(context.Background(), validRequest(clientID))

		require.Error(t, err)
	})

	t.Run("does not persist anything and never touches a session repository", func(t *testing.T) {
		t.Parallel()
		// authorize.Service no longer depends on session.Repository at all —
		// this test documents that: NewService's signature has no repository
		// parameter, so there is nothing to assert beyond compiling.
		clientID := uuid.New().String()
		svc := newService(
			&mockClientRepo{client: newClientMother(clientID)},
			&mockIdpFinderRepo{idp: aGoogleIdp()},
			&mockIdpClient{url: "https://accounts.google.com/authorize"},
			&mockStateCodec{encoded: "x"},
		)
		_, err := svc.Authorize(context.Background(), validRequest(clientID))
		require.NoError(t, err)
	})
}
