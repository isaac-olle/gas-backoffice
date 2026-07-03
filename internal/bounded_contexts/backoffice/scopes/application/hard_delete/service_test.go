package hard_delete_test

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
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/hard_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func TestHardDelete(t *testing.T) {
	t.Parallel()

	t.Run("hard deletes the found scope with correct id and publishes ScopeDeleted event", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		scope, scopeId := aScope(t)

		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(scope, nil)
		repo.EXPECT().HardDelete(mock.Anything, scopeId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ScopeDeletedType
			}))

		finder := domain_services.NewScopeFinder(repo)
		svc := hard_delete.NewService(finder, repo, bus)

		require.NoError(t, svc.HardDelete(context.Background(), aCommand(t, withScopeId(scopeId))))
		assert.Equal(t, scopeId.String(), scope.ToPrimitives().Id)
	})

	t.Run("returns error when finder returns error", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		_, scopeId := aScope(t)

		repo.EXPECT().SearchById(mock.Anything, scopeId).Return(nil, fmt.Errorf("db error"))

		finder := domain_services.NewScopeFinder(repo)
		svc := hard_delete.NewService(finder, repo, bus)
		err := svc.HardDelete(context.Background(), aCommand(t, withScopeId(scopeId)))

		assert.Error(t, err)
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *scopemocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrScopeNotSaved when hard delete fails",
			setupRepo: func(repo *scopemocks.MockRepository) {
				scope, _ := aScope(t)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(scope, nil)
				repo.EXPECT().
					HardDelete(mock.Anything, mock.Anything).
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

			finder := domain_services.NewScopeFinder(repo)
			svc := hard_delete.NewService(finder, repo, bus)
			err := svc.HardDelete(context.Background(), aCommand(t))

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
