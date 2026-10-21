package create_tenant_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/create_tenant"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

func newService(t *testing.T) (*create_tenant.Service, *tenantmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	repo := tenantmocks.NewMockRepository(t)
	bus := eventmocks.NewMockBus(t)
	return create_tenant.NewService(repo, bus), repo, bus
}

func TestCreateTenant(t *testing.T) {
	t.Parallel()

	t.Run("saves tenant and publishes UserCreated event", func(t *testing.T) {
		t.Parallel()
		svc, repo, bus := newService(t)

		var saved *tenants.Tenant
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(u *tenants.Tenant) bool {
				saved = u
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.UserCreatedType
			}))

		require.NoError(t, svc.Create(context.Background(), aCommand(t)))
		assert.Equal(t, "Tenant Name", saved.ToPrimitives().Name)
		assert.False(t, saved.IsSystemResource())
	})

	t.Run("saves system tenant when isSystem is set", func(t *testing.T) {
		t.Parallel()
		svc, repo, bus := newService(t)

		var saved *tenants.Tenant
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(u *tenants.Tenant) bool {
				saved = u
				return true
			})).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.Anything)

		require.NoError(t, svc.Create(context.Background(), aCommand(t, withIsSystem(true))))
		assert.True(t, saved.IsSystemResource())
	})

	errCases := []struct {
		desc      string
		cmd       *create_tenant.Command
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "empty name returns ErrUserNotCreated",
			cmd:      aCommand(t, withName("")),
			wantCode: tenants.ErrUserNotCreated,
		},
		{
			desc: "repository error returns ErrUserNotSaved",
			cmd:  aCommand(t),
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*tenants.Tenant")).
					Return(fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrUserNotSaved,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			svc, repo, bus := newService(t)
			_ = bus
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			err := svc.Create(context.Background(), tc.cmd)

			require.Error(t, err)
			var appErr *app_errors.AppError
			require.ErrorAs(t, err, &appErr)
			assert.Equal(t, tc.wantCode, appErr.Code)
		})
	}
}
