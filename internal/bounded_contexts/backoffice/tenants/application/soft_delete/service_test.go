package soft_delete_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/soft_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
)

func TestSoftDelete(t *testing.T) {
	t.Parallel()

	t.Run("soft deletes the found tenant and publishes UserDeleted event", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)

		userId := ids.NewUserId()
		tenant := aTenantWithId(t, userId)
		cmd := aCommand(t, withUserId(userId.String()))

		repo.EXPECT().SearchById(mock.Anything, userId).Return(tenant, nil)
		repo.EXPECT().SoftDelete(mock.Anything, userId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.UserDeletedType
			}))

		svc := soft_delete.NewService(domain_services.NewTenantFinder(repo), repo, bus)
		require.NoError(t, svc.SoftDelete(context.Background(), cmd))
		assert.Equal(t, userId.String(), tenant.ToPrimitives().Id)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *soft_delete.Command
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrUserNotFound when tenant does not exist",
			cmd:  func(t *testing.T) *soft_delete.Command { return aCommand(t) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), nil)
			},
			wantCode: tenants.ErrUserNotFound,
		},
		{
			desc: "returns ErrUserNotFound when repository returns error on find",
			cmd:  func(t *testing.T) *soft_delete.Command { return aCommand(t) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrUserNotFound,
		},
		{
			desc: "returns forbidden when caller is not self or owning admin",
			cmd: func(t *testing.T) *soft_delete.Command {
				return aCommand(t, withAuthor(&authors.UserAuthor{Id: ids.NewUserId(), Role: vo.TenantRole}))
			},
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(aTenant(t), nil)
			},
			wantCode: tenants.ErrInvalidOwnership,
		},
		{
			desc: "returns ErrUserNotSaved when soft delete fails",
			cmd: func(t *testing.T) *soft_delete.Command {
				userId := ids.NewUserId()
				return aCommand(t, withUserId(userId.String()))
			},
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, id ids.UserId) (*tenants.Tenant, error) {
						return aTenantWithId(t, id), nil
					})
				repo.EXPECT().
					SoftDelete(mock.Anything, mock.Anything).
					Return(fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrUserNotSaved,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := tenantmocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			svc := soft_delete.NewService(domain_services.NewTenantFinder(repo), repo, bus)
			err := svc.SoftDelete(context.Background(), tc.cmd(t))

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
