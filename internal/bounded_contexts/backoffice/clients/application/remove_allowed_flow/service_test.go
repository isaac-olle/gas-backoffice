package remove_allowed_flow_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/remove_allowed_flow"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

func TestRemoveAllowedFlow(t *testing.T) {
	t.Parallel()

	t.Run("removes flow from client and publishes ClientFlowRemoved event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		c, clientId := aClientWithFlow(t, "pkce")

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
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ClientFlowRemovedType
			}))

		cmd := remove_allowed_flow.NewCommand(clientId, "pkce")
		require.NoError(t, newService(t, repo, bus).Remove(context.Background(), cmd))

		assert.NotContains(t, saved.ToPrimitives().AllowedFlows, "pkce")
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *remove_allowed_flow.Command
		setupRepo func(repo *clientmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid flow name",
			cmd:      func(t *testing.T) *remove_allowed_flow.Command { return aCommand(t, withFlow("implicit")) },
			wantCode: vos.ErrInvalidFlow,
		},
		{
			desc: "returns not found when client does not exist",
			cmd:  func(t *testing.T) *remove_allowed_flow.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "propagates repository search error",
			cmd:  func(t *testing.T) *remove_allowed_flow.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "wraps repository save failure",
			cmd:  func(t *testing.T) *remove_allowed_flow.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.NewClient(ids2.NewClientId(), "valid-client", clients.NewOwnership(ids2.NewAdminId(), ids2.NewUserId()), nil, []string{"pkce"})
				c.PullDomainEvents()
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(fmt.Errorf("db error"))
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
			assertErrCode(t, newService(t, repo, bus).Remove(context.Background(), tc.cmd(t)), tc.wantCode)
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func newMocks(t *testing.T) (*clientmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	return clientmocks.NewMockRepository(t), eventmocks.NewMockBus(t)
}

func newService(t *testing.T, repo *clientmocks.MockRepository, bus *eventmocks.MockBus) *remove_allowed_flow.Service {
	t.Helper()
	return remove_allowed_flow.NewService(domain_services.NewClientFinder(repo), repo, bus)
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
