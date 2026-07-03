package create_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/create"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("saves scope with correct name and publishes ScopeCreated event", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		cmd := aCommand(t, withName("tenants:read"))

		var saved *scopes.Scope
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(s *scopes.Scope) bool {
				saved = s
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ScopeCreatedType
			}))

		require.NoError(t, create.NewService(repo, bus).Create(context.Background(), cmd))
		assert.Equal(t, "tenants:read", saved.ToPrimitives().Name)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *create.Command
		setupRepo func(repo *scopemocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects empty name with ErrScopeNotCreated",
			cmd:      func(t *testing.T) *create.Command { return aCommand(t, withName("")) },
			wantCode: scopes.ErrScopeNotCreated,
		},
		{
			desc: "wraps repository error as ErrScopeNotSaved",
			cmd:  func(t *testing.T) *create.Command { return aCommand(t) },
			setupRepo: func(repo *scopemocks.MockRepository) {
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

			svc := create.NewService(repo, bus)
			err := svc.Create(context.Background(), tc.cmd(t))

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
