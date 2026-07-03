package add_allowed_flow_test

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
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_allowed_flow"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"
)

func TestAddAllowedFlow(t *testing.T) {
	t.Parallel()

	t.Run("adds flow to client and publishes ClientFlowAdded event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		c, clientId := aClient(t)

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
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ClientFlowAddedType
			}))

		cmd := add_allowed_flow.NewCommand(clientId, "pkce")
		require.NoError(t, add_allowed_flow.NewService(repo, bus).Add(context.Background(), cmd))

		assert.Contains(t, saved.ToPrimitives().AllowedFlows, "pkce")
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *add_allowed_flow.Command
		setupRepo func(repo *clientmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid flow name",
			cmd:      func(t *testing.T) *add_allowed_flow.Command { return aCommand(t, withFlow("implicit")) },
			wantCode: vos.ErrInvalidFlow,
		},
		{
			desc: "propagates repository search failure",
			cmd:  func(t *testing.T) *add_allowed_flow.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
			},
			wantCode: "",
		},
		{
			desc: "wraps repository save failure",
			cmd:  func(t *testing.T) *add_allowed_flow.Command { return aCommand(t) },
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(
					func() *clients.Client {
						clientId := ids2.NewClientId()
						c, _ := clients.NewClient(clientId, "valid-client", clients.NewOwnership(ids2.NewAdminId(), ids2.NewTenantId()), nil, nil)
						c.PullDomainEvents()
						return c
					}(),
					nil,
				)
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
			err := add_allowed_flow.NewService(repo, bus).Add(context.Background(), tc.cmd(t))
			require.Error(t, err)
			if tc.wantCode != "" {
				assertErrCode(t, err, tc.wantCode)
			}
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func newMocks(t *testing.T) (*clientmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	return clientmocks.NewMockRepository(t), eventmocks.NewMockBus(t)
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
