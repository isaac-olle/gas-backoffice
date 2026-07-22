package check_system_resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/check_system_resource"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	usermocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func aUser(t *testing.T, isSystem bool) *users.User {
	t.Helper()
	id := ids.NewUserId().String()
	user, err := users.FromPrimitives(users.UserPrimitives{
		Id:       id,
		AdminId:  id,
		Name:     "Test User",
		Role:     "admin",
		IsSystem: isSystem,
	})
	require.NoError(t, err)
	return user
}

func TestCheckSystemResource(t *testing.T) {
	t.Parallel()

	userId := ids.NewUserId()
	qry := check_system_resource.NewQuery(userId)

	t.Run("returns forbidden when user is system resource", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.UserId")).Return(aUser(t, true), nil)

		svc := check_system_resource.NewService(domain_services.NewUserFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.Error(t, err)
	})

	t.Run("returns ok when user is not system resource", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.UserId")).Return(aUser(t, false), nil)

		svc := check_system_resource.NewService(domain_services.NewUserFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.NoError(t, err)
	})

	t.Run("propagates finder error", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.UserId")).Return(nil, fmt.Errorf("db error"))

		svc := check_system_resource.NewService(domain_services.NewUserFinder(repo))
		_, err := svc.Check(context.Background(), qry)
		assert.Error(t, err)
	})
}
