package check_system_resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/check_system_resource"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func aPR(t *testing.T, isSystem bool) *protected_resources.ProtectedResource {
	t.Helper()
	pr, err := protected_resources.FromPrimitives(protected_resources.ProtectedResourcePrimitives{
		Id:       ids.NewProtectedResourceId().String(),
		Name:     "Test Resource",
		TenantId: ids.NewTenantId().String(),
		AdminId:  ids.NewAdminId().String(),
		IsSystem: isSystem,
	})
	require.NoError(t, err)
	return pr
}

func TestCheckSystemResource(t *testing.T) {
	t.Parallel()

	prId := ids.NewProtectedResourceId()
	qry := check_system_resource.NewQuery(prId)

	t.Run("returns forbidden when protected resource is system resource", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ProtectedResourceId")).Return(aPR(t, true), nil)

		svc := check_system_resource.NewService(domain_services.NewProtectedResourceFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.Error(t, err)
	})

	t.Run("returns ok when protected resource is not system resource", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ProtectedResourceId")).Return(aPR(t, false), nil)

		svc := check_system_resource.NewService(domain_services.NewProtectedResourceFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.NoError(t, err)
	})

	t.Run("propagates finder error", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ProtectedResourceId")).Return(nil, fmt.Errorf("db error"))

		svc := check_system_resource.NewService(domain_services.NewProtectedResourceFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.Error(t, err)
	})
}
