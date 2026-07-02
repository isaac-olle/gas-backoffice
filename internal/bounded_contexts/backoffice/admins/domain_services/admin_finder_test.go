package domain_services_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins/domain_services"
	adminmocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// anAdmin builds a minimal valid Admin for test setup.
func anAdmin(t *testing.T) *admins.Admin {
	t.Helper()
	a, err := admins.FromPrimitives(admins.AdminPrimitives{
		Id:        ids.NewAdminId().String(),
		Name:      "Test",
		TenantIds: []string{},
	})
	require.NoError(t, err)
	return a
}

func TestAdminFinder(t *testing.T) {
	t.Parallel()

	t.Run("returns admin when repo finds it", func(t *testing.T) {
		t.Parallel()
		expected := anAdmin(t)
		repo := adminmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.AdminId")).Return(expected, nil)

		finder := domain_services.NewAdminFinder(repo)
		got, err := finder.Find(context.Background(), ids.NewAdminId())
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("propagates error from repository", func(t *testing.T) {
		t.Parallel()
		repo := adminmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.AdminId")).Return(nil, fmt.Errorf("db error"))

		finder := domain_services.NewAdminFinder(repo)
		_, err := finder.Find(context.Background(), ids.NewAdminId())
		assert.Error(t, err)
	})

	t.Run("returns not found error when repo returns nil", func(t *testing.T) {
		t.Parallel()
		repo := adminmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.AdminId")).Return(nil, nil)

		finder := domain_services.NewAdminFinder(repo)
		_, err := finder.Find(context.Background(), ids.NewAdminId())
		assert.Error(t, err)
	})
}
