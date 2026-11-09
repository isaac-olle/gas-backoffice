package domain_services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func aClient(t *testing.T) *clients.Client {
	t.Helper()
	c, err := clients.FromPrimitives(clients.ClientPrimitives{
		Id:             ids.NewClientId().String(),
		Name:           "Test Client",
		SubscriptionId: ids.NewSubscriptionId().String(),
	})
	require.NoError(t, err)
	return c
}

func TestClientFinder(t *testing.T) {
	t.Parallel()

	t.Run("returns client when repo finds it", func(t *testing.T) {
		t.Parallel()
		expected := aClient(t)
		repo := clientmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ClientId")).Return(expected, nil)

		finder := domain_services.NewClientFinder(repo)
		got, err := finder.Find(context.Background(), ids.NewClientId())
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("propagates error from repository", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ClientId")).Return(nil, fmt.Errorf("db error"))

		finder := domain_services.NewClientFinder(repo)
		_, err := finder.Find(context.Background(), ids.NewClientId())
		assert.Error(t, err)
	})

	t.Run("returns not found error when repo returns nil", func(t *testing.T) {
		t.Parallel()
		repo := clientmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ClientId")).Return(nil, nil)

		finder := domain_services.NewClientFinder(repo)
		_, err := finder.Find(context.Background(), ids.NewClientId())
		assert.Error(t, err)
	})
}
