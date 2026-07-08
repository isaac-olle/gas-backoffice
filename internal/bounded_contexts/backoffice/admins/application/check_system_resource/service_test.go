package check_system_resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/check_system_resource"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/domain_services"
	adminmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func anAdmin(t *testing.T, isSystem bool) *admins.Admin {
	t.Helper()
	a, err := admins.FromPrimitives(admins.AdminPrimitives{
		Id:       ids.NewAdminId().String(),
		Name:     "Test",
		UserIds:  []string{},
		IsSystem: isSystem,
	})
	require.NoError(t, err)
	return a
}

func TestCheckSystemResource(t *testing.T) {
	t.Parallel()

	adminId := ids.NewAdminId()
	qry := check_system_resource.NewQuery(adminId)

	t.Run("returns forbidden when admin is system resource", func(t *testing.T) {
		t.Parallel()
		repo := adminmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.AdminId")).Return(anAdmin(t, true), nil)

		svc := check_system_resource.NewService(domain_services.NewAdminFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.Error(t, err)
	})

	t.Run("returns ok when admin is not system resource", func(t *testing.T) {
		t.Parallel()
		repo := adminmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.AdminId")).Return(anAdmin(t, false), nil)

		svc := check_system_resource.NewService(domain_services.NewAdminFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.NoError(t, err)
	})

	t.Run("propagates finder error", func(t *testing.T) {
		t.Parallel()
		repo := adminmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.AdminId")).Return(nil, fmt.Errorf("db error"))

		svc := check_system_resource.NewService(domain_services.NewAdminFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.Error(t, err)
	})
}
