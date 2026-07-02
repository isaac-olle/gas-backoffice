package domain_services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys/domain_services"
	cakmocks "github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys/domain_services/mocks"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func aClientAssertionKey(t *testing.T) *client_assertion_keys.ClientAssertionKey {
	t.Helper()
	cak, err := client_assertion_keys.FromPrimitives(client_assertion_keys.ClientAssertionKeyPrimitives{
		Id:        ids.NewClientAssertionKeyId().String(),
		ClientId:  ids.NewClientId().String(),
		Algorithm: "RS256",
		PublicKey: []byte("stub-public-key"),
	})
	require.NoError(t, err)
	return cak
}

func TestClientAssertionKeyFinder_Find(t *testing.T) {
	t.Parallel()

	t.Run("returns key when repo finds it", func(t *testing.T) {
		t.Parallel()
		expected := aClientAssertionKey(t)
		repo := cakmocks.NewMockClientAssertionKeyRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ClientAssertionKeyId")).Return(expected, nil)

		finder := domain_services.NewClientAssertionKeyFinder(repo)
		got, err := finder.Find(context.Background(), ids.NewClientAssertionKeyId())
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("propagates error from repository", func(t *testing.T) {
		t.Parallel()
		repo := cakmocks.NewMockClientAssertionKeyRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ClientAssertionKeyId")).Return(nil, fmt.Errorf("db error"))

		finder := domain_services.NewClientAssertionKeyFinder(repo)
		_, err := finder.Find(context.Background(), ids.NewClientAssertionKeyId())
		assert.Error(t, err)
	})

	t.Run("returns not found error when repo returns nil", func(t *testing.T) {
		t.Parallel()
		repo := cakmocks.NewMockClientAssertionKeyRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.ClientAssertionKeyId")).Return(nil, nil)

		finder := domain_services.NewClientAssertionKeyFinder(repo)
		_, err := finder.Find(context.Background(), ids.NewClientAssertionKeyId())
		assert.Error(t, err)
	})
}

func TestClientAssertionKeyFinder_FindByClientId(t *testing.T) {
	t.Parallel()

	t.Run("returns key when repo finds it by client id", func(t *testing.T) {
		t.Parallel()
		expected := aClientAssertionKey(t)
		repo := cakmocks.NewMockClientAssertionKeyRepository(t)
		repo.EXPECT().FindByClientId(mock.Anything, mock.AnythingOfType("ids.ClientId")).Return(expected, nil)

		finder := domain_services.NewClientAssertionKeyFinder(repo)
		got, err := finder.FindByClientId(context.Background(), ids.NewClientId())
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("propagates error from repository", func(t *testing.T) {
		t.Parallel()
		repo := cakmocks.NewMockClientAssertionKeyRepository(t)
		repo.EXPECT().FindByClientId(mock.Anything, mock.AnythingOfType("ids.ClientId")).Return(nil, fmt.Errorf("db error"))

		finder := domain_services.NewClientAssertionKeyFinder(repo)
		_, err := finder.FindByClientId(context.Background(), ids.NewClientId())
		assert.Error(t, err)
	})

	t.Run("returns not found error when repo returns nil", func(t *testing.T) {
		t.Parallel()
		repo := cakmocks.NewMockClientAssertionKeyRepository(t)
		repo.EXPECT().FindByClientId(mock.Anything, mock.AnythingOfType("ids.ClientId")).Return(nil, nil)

		finder := domain_services.NewClientAssertionKeyFinder(repo)
		_, err := finder.FindByClientId(context.Background(), ids.NewClientId())
		assert.Error(t, err)
	})
}
