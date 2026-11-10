package remove_allowed_grant_type_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/remove_allowed_grant_type"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

func TestRemoveAllowedGrantType(t *testing.T) {
	t.Parallel()

	t.Run("removes grant type from client and publishes ClientGrantTypeRemoved event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		subscriptionId := ids2.NewSubscriptionId()
		c, clientId := aClientWithGrantType(t, subscriptionId, "client_credentials")
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
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ClientGrantTypeRemovedType
			}))

		cmd := remove_allowed_grant_type.NewCommand(clientId, "client_credentials", author)
		require.NoError(t, newService(t, repo, bus).Remove(context.Background(), cmd))

		assert.NotContains(t, saved.ToPrimitives().AllowedGrantTypes, "client_credentials")
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *remove_allowed_grant_type.Command
		setupRepo func(repo *clientmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid grant type name",
			cmd:      func(t *testing.T) *remove_allowed_grant_type.Command { return aCommand(t, withGrantType("implicit")) },
			wantCode: vos.ErrInvalidGrantType,
		},
		{
			desc: "returns not found when client does not exist",
			cmd:  func(t *testing.T) *remove_allowed_grant_type.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "propagates repository search error",
			cmd:  func(t *testing.T) *remove_allowed_grant_type.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns forbidden when author lacks subscription access",
			cmd: func(t *testing.T) *remove_allowed_grant_type.Command {
				unauthorized, err := vo.NewUserAuthorFromString(ids2.NewUserId().String(), ids2.NewUserId().String(), vo.SubscriptionUserRole)
				require.NoError(t, err)
				return aCommand(t, withAuthor(unauthorized))
			},
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := aClientWithGrantType(t, ids2.NewSubscriptionId(), "client_credentials")
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrInvalidOwnership,
		},
		{
			desc: "returns forbidden when client is system resource",
			cmd:  func(t *testing.T) *remove_allowed_grant_type.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, err := clients.NewClient(ids2.NewClientId(), "system-client", clients.NewOwnership(ids2.NewSubscriptionId()), nil, []string{"client_credentials"}, false, true)
				require.NoError(t, err)
				c.PullDomainEvents()
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrSystemResource,
		},
		{
			desc: "wraps repository save failure",
			cmd:  func(t *testing.T) *remove_allowed_grant_type.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := aClientWithGrantType(t, ids2.NewSubscriptionId(), "client_credentials")
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotSaved,
		},
		{
			desc: "returns conflict when repository reports concurrent modification",
			cmd:  func(t *testing.T) *remove_allowed_grant_type.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := aClientWithGrantType(t, ids2.NewSubscriptionId(), "client_credentials")
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
			assertErrCode(t, newService(t, repo, bus).Remove(context.Background(), tc.cmd(t)), tc.wantCode)
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func newMocks(t *testing.T) (*clientmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	return clientmocks.NewMockRepository(t), eventmocks.NewMockBus(t)
}

func newService(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus) *remove_allowed_grant_type.Service {
	t.Helper()
	finder := domain_services.NewClientFinder(repo)
	return remove_allowed_grant_type.NewService(finder, repo, bus, sharedDomainServices.NewSubscriptionAuthorizer())
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
