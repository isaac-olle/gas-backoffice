package soft_delete_test

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
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/soft_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
)

func TestSoftDelete(t *testing.T) {
	t.Parallel()

	t.Run("soft deletes the found tenant and publishes TenantDeleted event", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)

		tenantId := ids.NewTenantId()
		tenant := aTenantWithId(t, tenantId)
		cmd := aCommand(t, withTenantId(tenantId.String()))

		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)
		repo.EXPECT().SoftDelete(mock.Anything, tenantId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.TenantDeletedType
			}))

		svc := soft_delete.NewService(domain_services.NewTenantFinder(repo), repo, bus)
		require.NoError(t, svc.SoftDelete(context.Background(), cmd))
		assert.Equal(t, tenantId.String(), tenant.ToPrimitives().Id)
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrTenantNotFound when tenant does not exist",
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), nil)
			},
			wantCode: tenants.ErrTenantNotFound,
		},
		{
			desc: "returns ErrTenantNotFound when repository returns error on find",
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrTenantNotFound,
		},
		{
			desc: "returns ErrTenantNotSaved when soft delete fails",
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(aTenant(t), nil)
				repo.EXPECT().
					SoftDelete(mock.Anything, mock.Anything).
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

			svc := soft_delete.NewService(domain_services.NewTenantFinder(repo), repo, bus)
			err := svc.SoftDelete(context.Background(), aCommand(t))

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
