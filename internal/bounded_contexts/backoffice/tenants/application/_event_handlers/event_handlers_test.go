package _event_handlers

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
)

// ---- AddClientOnClientCreated ----

func TestAddClientOnClientCreated(t *testing.T) {
	t.Parallel()

	t.Run("adds client id to tenant, saves updated tenant, and publishes events", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		tenantId := ids2.NewTenantId()
		adminId := ids2.NewAdminId()
		clientId := ids2.NewClientId()
		tenant := aTenant(t, tenantId, adminId)

		var saved *tenants.Tenant
		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *tenants.Tenant) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
			return true // AddClientId does not record domain events; empty slice is expected
		}))

		handler := NewAddClientOnClientCreatedEventHandler(repo, bus)
		evt := backoffice.NewClientCreated(clientId.String(), "Client", adminId.String(), tenantId.String())

		require.NoError(t, handler.Handle(context.Background(), evt))

		assert.Contains(t, saved.ToPrimitives().ClientIds, clientId.String())
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrTenantNotFound when repository find fails",
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), fmt.Errorf("not found"))
			},
			wantCode: tenants.ErrTenantNotFound,
		},
		{
			desc: "returns ErrTenantNotSaved when save fails",
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return(aTenant(t, ids2.NewTenantId(), ids2.NewAdminId()), nil)
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*tenants.Tenant")).
					Return(fmt.Errorf("save error"))
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

			tenantId := ids2.NewTenantId()
			adminId := ids2.NewAdminId()
			clientId := ids2.NewClientId()

			handler := NewAddClientOnClientCreatedEventHandler(repo, bus)
			evt := backoffice.NewClientCreated(clientId.String(), "Client", adminId.String(), tenantId.String())

			err := handler.Handle(context.Background(), evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---- AddProtectedResourceOnProtectedResourceCreated ----

func TestAddProtectedResourceOnProtectedResourceCreated(t *testing.T) {
	t.Parallel()

	t.Run("adds protected resource id to tenant, saves updated tenant, and publishes events", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		tenantId := ids2.NewTenantId()
		adminId := ids2.NewAdminId()
		prId := ids2.NewProtectedResourceId()
		tenant := aTenant(t, tenantId, adminId)

		var saved *tenants.Tenant
		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *tenants.Tenant) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
			return true // AddProtectedResourceId does not record domain events; empty slice is expected
		}))

		handler := NewAddProtectedResourceOnProtectedResourceCreatedEventHandler(repo, bus)
		evt := backoffice.NewProtectedResourceCreated(prId.String(), "PR", tenantId.String())

		require.NoError(t, handler.Handle(context.Background(), evt))

		assert.Contains(t, saved.ToPrimitives().ProtectedResourceIds, prId.String())
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrTenantNotFound when repository find fails",
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), fmt.Errorf("not found"))
			},
			wantCode: tenants.ErrTenantNotFound,
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

			tenantId := ids2.NewTenantId()
			prId := ids2.NewProtectedResourceId()

			handler := NewAddProtectedResourceOnProtectedResourceCreatedEventHandler(repo, bus)
			evt := backoffice.NewProtectedResourceCreated(prId.String(), "PR", tenantId.String())

			err := handler.Handle(context.Background(), evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---- RemoveClientOnClientDeleted ----

func TestRemoveClientOnClientDeleted(t *testing.T) {
	t.Parallel()

	t.Run("removes client id from tenant, saves updated tenant, and publishes events", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		tenantId := ids2.NewTenantId()
		adminId := ids2.NewAdminId()
		clientId := ids2.NewClientId()

		prims := tenants.TenantPrimitives{
			Id:        tenantId.String(),
			Name:      "Test Tenant",
			AdminId:   adminId.String(),
			ClientIds: []string{clientId.String()},
		}
		tenant, err := tenants.FromPrimitives(prims)
		require.NoError(t, err)

		var saved *tenants.Tenant
		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *tenants.Tenant) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
			return true // RemoveClientId does not record domain events; empty slice is expected
		}))

		handler := NewRemoveClientOnClientDeletedEventHandler(repo, bus)
		evt := backoffice.NewClientDeleted(clientId.String(), tenantId.String())

		require.NoError(t, handler.Handle(context.Background(), evt))
		assert.Empty(t, saved.ToPrimitives().ClientIds)
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrTenantNotFound when repository find fails",
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrTenantNotFound,
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

			tenantId := ids2.NewTenantId()
			clientId := ids2.NewClientId()

			handler := NewRemoveClientOnClientDeletedEventHandler(repo, bus)
			evt := backoffice.NewClientDeleted(clientId.String(), tenantId.String())

			err := handler.Handle(context.Background(), evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---- RemoveProtectedResourceOnProtectedResourceDeleted ----

func TestRemoveProtectedResourceOnProtectedResourceDeleted(t *testing.T) {
	t.Parallel()

	t.Run("removes protected resource id from tenant, saves updated tenant, and publishes events", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		tenantId := ids2.NewTenantId()
		adminId := ids2.NewAdminId()
		prId := ids2.NewProtectedResourceId()

		prims := tenants.TenantPrimitives{
			Id:                   tenantId.String(),
			Name:                 "Test Tenant",
			AdminId:              adminId.String(),
			ProtectedResourceIds: []string{prId.String()},
		}
		tenant, err := tenants.FromPrimitives(prims)
		require.NoError(t, err)

		var saved *tenants.Tenant
		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(ten *tenants.Tenant) bool {
				saved = ten
				return true
			})).Return(nil)
		bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
			return true // RemoveProtectedResourceId does not record domain events; empty slice is expected
		}))

		handler := NewRemoveProtectedResourceOnProtectedResourceDeletedEventHandler(repo, bus)
		evt := backoffice.NewProtectedResourceDeleted(prId.String(), tenantId.String())

		require.NoError(t, handler.Handle(context.Background(), evt))
		assert.Empty(t, saved.ToPrimitives().ProtectedResourceIds)
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *tenantmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns ErrTenantNotFound when repository find fails",
			setupRepo: func(repo *tenantmocks.MockRepository) {
				repo.EXPECT().
					SearchById(mock.Anything, mock.Anything).
					Return((*tenants.Tenant)(nil), fmt.Errorf("db error"))
			},
			wantCode: tenants.ErrTenantNotFound,
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

			tenantId := ids2.NewTenantId()
			prId := ids2.NewProtectedResourceId()

			handler := NewRemoveProtectedResourceOnProtectedResourceDeletedEventHandler(repo, bus)
			evt := backoffice.NewProtectedResourceDeleted(prId.String(), tenantId.String())

			err := handler.Handle(context.Background(), evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// --- helpers ---

func aTenant(t *testing.T, tenantId ids2.TenantId, adminId ids2.AdminId) *tenants.Tenant {
	t.Helper()
	tenant, err := tenants.NewTenant(tenantId, "Test Tenant", tenants.NewOwnership(adminId))
	require.NoError(t, err)
	tenant.PullDomainEvents()
	return tenant
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
