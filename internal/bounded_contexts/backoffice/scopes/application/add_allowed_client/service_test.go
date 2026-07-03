package add_allowed_client_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func TestAddAllowedClient(t *testing.T) {
	t.Parallel()

	t.Run("adds client to scope, saves updated scope, and publishes ScopeClientAdded event", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		scope, scopeId := aScope(t)
		clientId := ids.NewClientId()

		var saved *scopes.Scope
		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(s *scopes.Scope) bool {
				saved = s
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ScopeClientAddedType
			}))

		svc := add_allowed_client.NewService(repo, bus)
		require.NoError(t, svc.AddAllowedClient(context.Background(),
			aCommand(t, withScopeId(scopeId.String()), withClientId(clientId.String()))))

		assert.Contains(t, saved.ToPrimitives().AllowedClientIds, clientId.String())
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *add_allowed_client.Command
		setupRepo func(repo *scopemocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrScopeNotFound for invalid scope id",
			cmd:      func(t *testing.T) *add_allowed_client.Command { return aCommand(t, withScopeId("not-a-uuid")) },
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc:     "returns ErrScopeNotFound for invalid client id",
			cmd:      func(t *testing.T) *add_allowed_client.Command { return aCommand(t, withClientId("not-a-uuid")) },
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc: "returns ErrScopeNotFound when scope does not exist",
			cmd:  func(t *testing.T) *add_allowed_client.Command { return aCommand(t) },
			setupRepo: func(repo *scopemocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*scopes.Scope)(nil), fmt.Errorf("not found"))
			},
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc: "returns ErrScopeNotSaved when save fails",
			cmd:  func(t *testing.T) *add_allowed_client.Command { return aCommand(t) },
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope, _ := aScope(t)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(scope, nil)
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*scopes.Scope")).
					Return(fmt.Errorf("db error"))
			},
			wantCode: scopes.ErrScopeNotSaved,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := scopemocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			svc := add_allowed_client.NewService(repo, bus)
			err := svc.AddAllowedClient(context.Background(), tc.cmd(t))

			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// --- helpers ---

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
