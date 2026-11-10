package add_redirect_uri_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_redirect_uri"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"
)

func TestAddRedirectUri(t *testing.T) {
	t.Parallel()

	t.Run("adds redirect uri to client and publishes ClientRedirectUriAdded event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		subscriptionId := ids2.NewSubscriptionId()
		c, clientId := aClient(t, subscriptionId)
		author, err := vo.NewUserAuthorWithMemberships(ids2.NewUserId().String(), ids2.NewUserId().String(), vo.SubscriptionUserRole,
			[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
		require.NoError(t, err)

		var saved *clients.Client
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(cl *clients.Client) bool {
				saved = cl
				return true
			})).
			Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ClientRedirectUriAddedType
			}))

		cmd := add_redirect_uri.NewCommand(clientId, "https://example.com/callback", author)
		require.NoError(t, newService(t, repo, bus).Add(context.Background(), cmd))

		assert.Contains(t, saved.ToPrimitives().RedirectUris, "https://example.com/callback")
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *add_redirect_uri.Command
		setupRepo func(repo *clientmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid redirect uri scheme",
			cmd:      func(t *testing.T) *add_redirect_uri.Command { return aCommand(t, withRedirectUri("ftp://bad.com")) },
			wantCode: vos.ErrInvalidRedirectUri,
		},
		{
			desc: "propagates repository search failure",
			cmd:  func(t *testing.T) *add_redirect_uri.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns not found when client does not exist",
			cmd:  func(t *testing.T) *add_redirect_uri.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns forbidden when author lacks subscription access",
			cmd: func(t *testing.T) *add_redirect_uri.Command {
				unauthorized, err := vo.NewUserAuthorFromString(ids2.NewUserId().String(), ids2.NewUserId().String(), vo.SubscriptionUserRole)
				require.NoError(t, err)
				return aCommand(t, withAuthor(unauthorized))
			},
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := aClient(t, ids2.NewSubscriptionId())
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrInvalidOwnership,
		},
		{
			desc: "returns forbidden when client is system resource",
			cmd:  func(t *testing.T) *add_redirect_uri.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, err := clients.NewClient(ids2.NewClientId(), "system-client", clients.NewOwnership(ids2.NewSubscriptionId()), nil, nil, false, true)
				require.NoError(t, err)
				c.PullDomainEvents()
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrSystemResource,
		},
		{
			desc: "wraps repository save failure",
			cmd:  func(t *testing.T) *add_redirect_uri.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := aClient(t, ids2.NewSubscriptionId())
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotSaved,
		},
		{
			desc: "returns conflict when repository reports concurrent modification",
			cmd:  func(t *testing.T) *add_redirect_uri.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := aClient(t, ids2.NewSubscriptionId())
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(persistance.ErrConcurrentModification)
			},
			wantCode: clients.ErrClientConcurrentModification,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo, bus := newMocks(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			assertErrCode(t, newService(t, repo, bus).Add(context.Background(), tc.cmd(t)), tc.wantCode)
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func newMocks(t *testing.T) (*clientmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	return clientmocks.NewMockRepository(t), eventmocks.NewMockBus(t)
}

func newService(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus) *add_redirect_uri.Service {
	t.Helper()
	finder := domain_services.NewClientFinder(repo)
	return add_redirect_uri.NewService(finder, repo, bus, sharedDomainServices.NewSubscriptionAuthorizer())
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
