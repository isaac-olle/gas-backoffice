package get_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/domain_services"
	idpmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("returns provider when found and enabled", func(t *testing.T) {
		t.Parallel()
		repo := idpmocks.NewMockRepository(t)
		id := ids.NewIdentityProviderId()
		repo.EXPECT().FindById(mock.Anything, id).Return(anIdentityProvider(t, id, true), nil)

		svc := get.NewService(domain_services.NewIdentityProviderFinder(repo))
		result, err := svc.Get(context.Background(), get.NewQuery(id))
		require.NoError(t, err)

		idp, ok := result.(*identity_providers.IdentityProvider)
		require.True(t, ok)
		assert.True(t, idp.IsEnabled())
	})

	t.Run("returns error when disabled", func(t *testing.T) {
		t.Parallel()
		repo := idpmocks.NewMockRepository(t)
		id := ids.NewIdentityProviderId()
		repo.EXPECT().FindById(mock.Anything, id).Return(anIdentityProvider(t, id, false), nil)

		svc := get.NewService(domain_services.NewIdentityProviderFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(id))
		assert.Error(t, err)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		t.Parallel()
		repo := idpmocks.NewMockRepository(t)
		id := ids.NewIdentityProviderId()
		repo.EXPECT().FindById(mock.Anything, id).Return(nil, nil)

		svc := get.NewService(domain_services.NewIdentityProviderFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(id))
		assert.Error(t, err)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := idpmocks.NewMockRepository(t)
		id := ids.NewIdentityProviderId()
		repo.EXPECT().FindById(mock.Anything, id).Return(nil, fmt.Errorf("db error"))

		svc := get.NewService(domain_services.NewIdentityProviderFinder(repo))
		_, err := svc.Get(context.Background(), get.NewQuery(id))
		assert.Error(t, err)
	})
}
