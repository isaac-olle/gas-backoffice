package get_test

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/application/get"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type getProviderRepository struct {
	idp *identity_providers.IdentityProvider
	err error
}

func (r *getProviderRepository) Save(context.Context, *identity_providers.IdentityProvider) error {
	return nil
}
func (r *getProviderRepository) FindById(context.Context, ids.IdentityProviderId) (*identity_providers.IdentityProvider, error) {
	return r.idp, r.err
}
func (r *getProviderRepository) FindByReferenceName(context.Context, string) (*identity_providers.IdentityProvider, error) {
	return nil, nil
}

func TestGet(t *testing.T) {
	id := ids.NewIdentityProviderId()
	tests := []struct {
		name     string
		provider *identity_providers.IdentityProvider
		repoErr  error
		wantErr  bool
	}{
		{name: "returns an enabled provider", provider: anIdentityProvider(t, id, true)},
		{name: "rejects disabled providers", provider: anIdentityProvider(t, id, false), wantErr: true},
		{name: "returns not found for missing provider", wantErr: true},
		{name: "wraps repository failures", repoErr: errors.New("query failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &getProviderRepository{idp: tt.provider, err: tt.repoErr}
			result, err := get.NewService(domain_services.NewIdentityProviderFinder(repo)).Get(context.Background(), get.NewQuery(id))
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				require.Same(t, tt.provider, result.(*identity_providers.IdentityProvider))
			}
		})
	}
}
