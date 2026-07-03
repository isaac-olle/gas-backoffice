package get_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
	tenantmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns tenant when found", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		tenantId := ids.NewTenantId()
		tenant := aTenantWithId(t, tenantId)

		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(tenant, nil)

		svc := get.NewService(domain_services.NewTenantFinder(repo))
		result, err := svc.Get(context.Background(), get.NewQuery(tenantId))
		require.NoError(t, err)

		ten, ok := result.(*tenants.Tenant)
		require.True(t, ok)
		assert.Equal(t, tenantId.String(), ten.ToPrimitives().Id)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		tenantId := ids.NewTenantId()

		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(nil, fmt.Errorf("db error"))

		svc := get.NewService(domain_services.NewTenantFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(tenantId))
		assert.Error(t, err)
	})

	t.Run("returns error when tenant not found", func(t *testing.T) {
		t.Parallel()
		repo := tenantmocks.NewMockRepository(t)
		tenantId := ids.NewTenantId()

		repo.EXPECT().SearchById(mock.Anything, tenantId).Return(nil, nil)

		svc := get.NewService(domain_services.NewTenantFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(tenantId))
		assert.Error(t, err)
	})
}
