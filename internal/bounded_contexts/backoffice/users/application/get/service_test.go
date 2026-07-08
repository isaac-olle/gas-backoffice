package get_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	usermocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns user when found", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		userId := ids.NewUserId()
		user := aUserWithId(t, userId)

		repo.EXPECT().SearchById(mock.Anything, userId).Return(user, nil)

		svc := get.NewService(domain_services.NewUserFinder(repo))
		result, err := svc.Get(context.Background(), get.NewQuery(userId))
		require.NoError(t, err)

		ten, ok := result.(*users.User)
		require.True(t, ok)
		assert.Equal(t, userId.String(), ten.ToPrimitives().Id)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		userId := ids.NewUserId()

		repo.EXPECT().SearchById(mock.Anything, userId).Return(nil, fmt.Errorf("db error"))

		svc := get.NewService(domain_services.NewUserFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(userId))
		assert.Error(t, err)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		t.Parallel()
		repo := usermocks.NewMockRepository(t)
		userId := ids.NewUserId()

		repo.EXPECT().SearchById(mock.Anything, userId).Return(nil, nil)

		svc := get.NewService(domain_services.NewUserFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(userId))
		assert.Error(t, err)
	})
}
