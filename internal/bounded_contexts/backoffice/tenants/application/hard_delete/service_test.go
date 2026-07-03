package hard_delete_test

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

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/hard_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
)

func TestHardDelete(t *testing.T) {
	t.Parallel()

	t.Run("hard deletes the found tenant and publishes TenantDeleted event", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)

		tenantId := ids.NewTenantId()
		tenant := aTenantWithId(t, tenantId)
		cmd := aCommand(t, withTenantId(tenantId.String()))

		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)
		repo.EXPECT().HardDelete(mock.Anything, tenantId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.TenantDeletedType
			}))

		svc := hard_delete.NewService(domain_services.NewTenantFinder(repo), repo, bus)
		require.NoError(t, svc.HardDelete(context.Background(), cmd))
		assert.Equal(t, tenantId.String(), tenant.ToPrimitives().Id)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *hard_delete.Command
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrTenantNotFound for invalid tenant id format",
			cmd:      func(t *testing.T) *hard_delete.Command { return aCommand(t, withTenantId("bad-id")) },
			wantCode: tenants.ErrTenantNotFound,
		},
		{
			desc: "returns ErrTenantNotFound when tenant does not exist",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), nil)
			},
			wantCode: tenants.ErrTenantNotFound,
		},
		{
			desc: "returns ErrTenantNotFound when repository find returns error",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrTenantNotFound,
		},
		{
			desc: "returns ErrTenantNotSaved when hard delete fails",
			cmd:  func(t *testing.T) *hard_delete.Command { return aCommand(t) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(aTenant(t), nil)
				repo.EXPECT().
					HardDelete(mock.Anything, mock.Anything).
					Return(fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrTenantNotSaved,
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

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
