package soft_delete_test

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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/soft_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
)

func TestSoftDelete(t *testing.T) {
	t.Parallel()

	t.Run("soft deletes client and publishes ClientDeleted event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		subscriptionId := ids2.NewSubscriptionId()
		c, clientId := aClient(t, subscriptionId)
		author, err := vo.NewUserAuthorWithMemberships(ids2.NewUserId().String(), ids2.NewUserId().String(), vo.SubscriptionUserRole,
			[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
		require.NoError(t, err)

		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)
		repo.EXPECT().SoftDelete(mock.Anything, clientId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ClientDeletedType
			}))

		err = newService(t, repo, bus).SoftDelete(context.Background(), soft_delete.NewCommand(clientId.String(), author))
		require.NoError(t, err)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *soft_delete.Command
		setupRepo func(repo *clientmocks.MockRepository) // nil = SearchById no s'hauria de cridar
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid client id format",
			cmd:      func(t *testing.T) *soft_delete.Command { return aCommand(t, withClientId("not-a-uuid")) },
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns not found when client does not exist",
			cmd:  func(t *testing.T) *soft_delete.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "propagates repository search error",
			cmd:  func(t *testing.T) *soft_delete.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns forbidden when author lacks subscription access",
			cmd: func(t *testing.T) *soft_delete.Command {
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
			cmd:  func(t *testing.T) *soft_delete.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, err := clients.NewClient(ids2.NewClientId(), "system-client", clients.NewOwnership(ids2.NewSubscriptionId()), nil, nil, false, true)
				require.NoError(t, err)
				c.PullDomainEvents()
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: clients.ErrSystemResource,
		},
		{
			desc: "wraps repository soft delete failure",
			cmd:  func(t *testing.T) *soft_delete.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := aClient(t, ids2.NewSubscriptionId())
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
				repo.EXPECT().SoftDelete(mock.Anything, mock.AnythingOfType("ids.ClientId")).Return(fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotSaved,
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
			assertErrCode(t, newService(t, repo, bus).SoftDelete(context.Background(), tc.cmd(t)), tc.wantCode)
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func newMocks(t *testing.T) (*clientmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	return clientmocks.NewMockRepository(t), eventmocks.NewMockBus(t)
}

func newService(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus) *soft_delete.Service {
	t.Helper()
	finder := domain_services.NewClientFinder(repo)
	return soft_delete.NewService(finder, repo, bus, sharedDomainServices.NewSubscriptionAuthorizer())
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
