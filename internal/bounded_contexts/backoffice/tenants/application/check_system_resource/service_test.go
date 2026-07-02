package check_system_resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/application/check_system_resource"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/domain_services"
	tenantmocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/mocks"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func aTenant(t *testing.T, isSystem bool) *tenants.Tenant {
	t.Helper()
	tenant, err := tenants.FromPrimitives(tenants.TenantPrimitives{
		Id:       ids.NewTenantId().String(),
		Name:     "Test Tenant",
		AdminId:  ids.NewAdminId().String(),
		IsSystem: isSystem,
	})
	require.NoError(t, err)
	return tenant
}

func TestCheckSystemResource(t *testing.T) {
	t.Parallel()

	tenantId := ids.NewTenantId()
	qry := check_system_resource.NewQuery(tenantId)

	t.Run("returns forbidden when tenant is system resource", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.TenantId")).Return(aTenant(t, true), nil)

		svc := check_system_resource.NewService(domain_services.NewTenantFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.Error(t, err)
	})

	t.Run("returns ok when tenant is not system resource", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.TenantId")).Return(aTenant(t, false), nil)

		svc := check_system_resource.NewService(domain_services.NewTenantFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.NoError(t, err)
	})

	t.Run("propagates finder error", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.TenantId")).Return(nil, fmt.Errorf("db error"))

		svc := check_system_resource.NewService(domain_services.NewTenantFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.Error(t, err)
	})
}
