package domain_services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type identityProviderFinderRepository struct {
	byID, byName   *identity_providers.IdentityProvider
	idErr, nameErr error
}

func (r *identityProviderFinderRepository) Save(context.Context, *identity_providers.IdentityProvider) error {
	return nil
}
func (r *identityProviderFinderRepository) FindById(context.Context, ids.IdentityProviderId) (*identity_providers.IdentityProvider, error) {
	return r.byID, r.idErr
}
func (r *identityProviderFinderRepository) FindByReferenceName(context.Context, string) (*identity_providers.IdentityProvider, error) {
	return r.byName, r.nameErr
}

func finderProvider(t *testing.T, id ids.IdentityProviderId, enabled bool) *identity_providers.IdentityProvider {
	t.Helper()
	idp, err := identity_providers.FromPrimitives(identity_providers.IdentityProviderPrimitives{Id: id.String(), ReferenceName: "google", Name: "Google", Issuer: "https://issuer.example", ClientId: "client", ClientSecret: "secret", AuthorizationEndpoint: "https://issuer.example/authorize", TokenEndpoint: "https://issuer.example/token", JwksUri: "https://issuer.example/jwks", Scopes: []string{"openid"}, BackofficeRedirectUri: "https://app.example/callback", Enabled: enabled})
	require.NoError(t, err)
	return idp
}

func TestIdentityProviderFinderFindEnabledByReferenceName(t *testing.T) {
	tests := []struct {
		name    string
		idp     *identity_providers.IdentityProvider
		repoErr error
		wantErr bool
	}{
		{name: "returns an enabled provider", idp: finderProvider(t, ids.NewIdentityProviderId(), true)},
		{name: "returns not found for nil result", wantErr: true},
		{name: "rejects disabled providers", idp: finderProvider(t, ids.NewIdentityProviderId(), false), wantErr: true},
		{name: "wraps repository errors", repoErr: errors.New("database unavailable"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &identityProviderFinderRepository{byName: tt.idp, nameErr: tt.repoErr}
			got, err := domain_services.NewIdentityProviderFinder(repo).FindEnabledByReferenceName(context.Background(), "google")
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				require.Same(t, tt.idp, got)
			}
		})
	}
}

func TestIdentityProviderFinderFindEnabled(t *testing.T) {
	id := ids.NewIdentityProviderId()
	tests := []struct {
		name    string
		idp     *identity_providers.IdentityProvider
		repoErr error
		wantErr bool
	}{
		{name: "returns an enabled provider", idp: finderProvider(t, id, true)},
		{name: "returns not found for nil result", wantErr: true},
		{name: "rejects disabled providers", idp: finderProvider(t, id, false), wantErr: true},
		{name: "wraps repository errors", repoErr: errors.New("database unavailable"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &identityProviderFinderRepository{byID: tt.idp, idErr: tt.repoErr}
			got, err := domain_services.NewIdentityProviderFinder(repo).FindEnabled(context.Background(), id)
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				require.Same(t, tt.idp, got)
			}
		})
	}
}
