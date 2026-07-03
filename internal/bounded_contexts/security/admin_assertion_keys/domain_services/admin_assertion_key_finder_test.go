package domain_services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys/domain_services"
	aakmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func anAdminAssertionKey(t *testing.T) *admin_assertion_keys.AdminAssertionKey {
	t.Helper()
	aak, err := admin_assertion_keys.FromPrimitives(admin_assertion_keys.AdminAssertionKeyPrimitives{
		Id:        ids.NewAdminAssertionKeyId().String(),
		AdminId:   ids.NewAdminId().String(),
		Algorithm: "RS256",
		PublicKey: []byte("stub-public-key"),
	})
	require.NoError(t, err)
	return aak
}

func TestAdminAssertionKeyFinder(t *testing.T) {
	t.Parallel()

	t.Run("returns key when repo finds it", func(t *testing.T) {
		t.Parallel()
		expected := anAdminAssertionKey(t)
		repo := aakmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.AdminAssertionKeyId")).Return(expected, nil)

		finder := domain_services.NewAdminAssertionKeyFinder(repo)
		got, err := finder.Find(context.Background(), ids.NewAdminAssertionKeyId())
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("propagates error from repository", func(t *testing.T) {
		t.Parallel()
		repo := aakmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.AdminAssertionKeyId")).Return(nil, fmt.Errorf("db error"))

		finder := domain_services.NewAdminAssertionKeyFinder(repo)
		_, err := finder.Find(context.Background(), ids.NewAdminAssertionKeyId())
		assert.Error(t, err)
	})

	t.Run("returns not found error when repo returns nil", func(t *testing.T) {
		t.Parallel()
		repo := aakmocks.NewMockRepository(t)
		repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.AdminAssertionKeyId")).Return(nil, nil)

		finder := domain_services.NewAdminAssertionKeyFinder(repo)
		_, err := finder.Find(context.Background(), ids.NewAdminAssertionKeyId())
		assert.Error(t, err)
	})
}
