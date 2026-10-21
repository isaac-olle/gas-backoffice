package hard_delete_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/hard_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
)

func TestHardDelete(t *testing.T) {
	t.Parallel()

	t.Run("hard deletes the found tenant and publishes UserDeleted event", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)

		userId := ids.NewUserId()
		tenant := aTenantWithId(t, userId)
		cmd := aCommand(t, withUserId(userId.String()), withAuthor(selfAuthor(userId)))

		repo.EXPECT().SearchById(mock.Anything, userId).Return(tenant, nil)
		repo.EXPECT().HardDelete(mock.Anything, userId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.UserDeletedType
			}))

		svc := hard_delete.NewService(domain_services.NewTenantFinder(repo), repo, bus)
		require.NoError(t, svc.HardDelete(context.Background(), cmd))
		assert.Equal(t, userId.String(), tenant.ToPrimitives().Id)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *hard_delete.Command
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrUserNotFound for invalid user id format",
			cmd:      func(t *testing.T) *hard_delete.Command { return aCommand(t, withUserId("bad-id")) },
			wantCode: tenants.ErrUserNotFound,
		},
		{
			desc: "returns ErrUserNotFound when tenant does not exist",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), nil)
			},
			wantCode: tenants.ErrUserNotFound,
		},
		{
			desc: "returns ErrUserNotFound when repository find returns error",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrUserNotFound,
		},
		{
			desc: "returns forbidden when caller has no access to the tenant",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(aTenant(t), nil)
			},
			wantCode: tenants.ErrInvalidOwnership,
		},
		{
			desc: "returns forbidden when tenant is system resource",
			cmd: func(t *testing.T) *hard_delete.Command {
				return aCommand(t, withUserId(systemResourceUserId.String()), withAuthor(selfAuthor(systemResourceUserId)))
			},
			setupRepo: func(repo *tenantmocks.MockRepository) {
				tenant, err := tenants.FromPrimitives(tenants.TenantPrimitives{
					Id:       systemResourceUserId.String(),
					AdminId:  ids.NewUserId().String(),
					Name:     "system-tenant",
					IsSystem: true,
				})
				require.NoError(t, err)
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(tenant, nil)
			},
			wantCode: tenants.ErrSystemResource,
		},
		{
			desc: "returns ErrUserNotSaved when hard delete fails",
			cmd: func(t *testing.T) *hard_delete.Command {
				userId := ids.NewUserId()
				return aCommand(t, withUserId(userId.String()), withAuthor(selfAuthor(userId)))
			},
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, id ids.UserId) (*tenants.Tenant, error) {
						return aTenantWithId(t, id), nil
					})
				repo.EXPECT().
					HardDelete(mock.Anything, mock.Anything).
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

			svc := hard_delete.NewService(domain_services.NewTenantFinder(repo), repo, bus)
			err := svc.HardDelete(context.Background(), tc.cmd(t))

			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// --- helpers ---

var systemResourceUserId = ids.NewUserId()

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
