package create_test

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

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/create"
	scopeDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func newService(scopeRepo *scopemocks.MockRepository, bus *eventmocks.MockBus) *create.Service {
	duplicateNameChecker := scopeDomainServices.NewDuplicateNameChecker(scopeRepo)
	return create.NewService(duplicateNameChecker, scopeRepo, bus)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("saves scope with correct name and publishes ScopeCreated event", func(t *testing.T) {
		t.Parallel()
		scopeRepo := scopemocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)

		prId := ids2.NewProtectedResourceId()
		userId := ids2.NewUserId()

		cmd := aCommand(t, withName("tenants:read"), withProtectedResourceId(prId), withUserId(userId))

		scopeRepo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, nil)

		var saved *scopes.Scope
		scopeRepo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(s *scopes.Scope) bool {
				saved = s
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ScopeCreatedType
			}))

		require.NoError(t, newService(scopeRepo, bus).Create(context.Background(), cmd))
		assert.Equal(t, "tenants:read", saved.ToPrimitives().Name)
	})

	prId := ids2.NewProtectedResourceId()
	userId := ids2.NewUserId()

	errCases := []struct {
		desc       string
		cmd        func(t *testing.T) *create.Command
		setupScope func(scopeRepo *scopemocks.MockRepository)
		wantCode   app_errors.ErrorCode
	}{
		{
			desc: "rejects empty name with ErrScopeNotCreated",
			cmd: func(t *testing.T) *create.Command {
				return aCommand(t, withName(""), withProtectedResourceId(prId), withUserId(userId))
			},
			wantCode: scopes.ErrScopeNotCreated,
		},
		{
			desc: "wraps repository error as ErrScopeNotSaved",
			cmd: func(t *testing.T) *create.Command {
				return aCommand(t, withProtectedResourceId(prId), withUserId(userId))
			},
			setupScope: func(scopeRepo *scopemocks.MockRepository) {
				scopeRepo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, nil)
				scopeRepo.EXPECT().
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
			scopeRepo := scopemocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			if tc.setupScope != nil {
				tc.setupScope(scopeRepo)
			}

			svc := newService(scopeRepo, bus)
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
