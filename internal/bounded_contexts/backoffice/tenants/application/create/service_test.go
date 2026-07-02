package create_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/buses/event/mocks"
	backofficeEvents "github.com/isaac-olle/gas/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/application/create"
	tenantmocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("saves tenant with correct data and publishes TenantCreated event", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		cmd := aCommand(t, withName("my-tenant"))

		var saved *tenants.Tenant
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *tenants.Tenant) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.TenantCreatedType
			}))

		require.NoError(t, create.NewService(repo, bus).Create(context.Background(), cmd))
		assert.Equal(t, "my-tenant", saved.ToPrimitives().Name)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *create.Command
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects empty name with ErrTenantNotCreated",
			cmd:      func(t *testing.T) *create.Command { return aCommand(t, withName("")) },
			wantCode: tenants.ErrTenantNotCreated,
		},
		{
			desc:     "rejects reserved name with ErrTenantNotCreated",
			cmd:      func(t *testing.T) *create.Command { return aCommand(t, withName("Gas-Tenant")) },
			wantCode: tenants.ErrTenantNotCreated,
		},
		{
			desc: "wraps repository error as ErrTenantNotSaved",
			cmd:  func(t *testing.T) *create.Command { return aCommand(t) },
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*tenants.Tenant")).
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
