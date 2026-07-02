package check_system_resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/application/check_system_resource"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/domain_services"
	scopemocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/mocks"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func aScope(t *testing.T, isSystem bool) *scopes.Scope {
	t.Helper()
	scope, err := scopes.FromPrimitives(scopes.ScopePrimitives{
		Id:                  ids.NewScopeId().String(),
		Name:                "read:users",
		ProtectedResourceId: ids.NewProtectedResourceId().String(),
		TenantId:            ids.NewTenantId().String(),
		AdminId:             ids.NewAdminId().String(),
		IsSystem:            isSystem,
	})
	require.NoError(t, err)
	return scope
}

func TestCheckSystemResource(t *testing.T) {
	t.Parallel()

	scopeId := ids.NewScopeId()
	qry := check_system_resource.NewQuery(scopeId)

	t.Run("returns forbidden when scope is system resource", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ScopeId")).Return(aScope(t, true), nil)

		svc := check_system_resource.NewService(domain_services.NewScopeFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.Error(t, err)
	})

	t.Run("returns ok when scope is not system resource", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ScopeId")).Return(aScope(t, false), nil)

		svc := check_system_resource.NewService(domain_services.NewScopeFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.NoError(t, err)
	})

	t.Run("propagates finder error", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ScopeId")).Return(nil, fmt.Errorf("db error"))

		svc := check_system_resource.NewService(domain_services.NewScopeFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.Error(t, err)
	})
}
