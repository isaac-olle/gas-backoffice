package soft_delete_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

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
		c, clientId := aClient(t)

		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)
		repo.EXPECT().SoftDelete(mock.Anything, clientId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ClientDeletedType
			}))

		err := newService(t, repo, bus).SoftDelete(context.Background(), soft_delete.NewCommand(clientId.String()))
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
			desc: "wraps repository soft delete failure",
			cmd:  func(t *testing.T) *soft_delete.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.NewClient(ids2.NewClientId(), "valid-client", clients.NewOwnership(ids2.NewUserId()), nil, nil)
				c.PullDomainEvents()
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
	return soft_delete.NewService(domain_services.NewClientFinder(repo), repo, bus)
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
