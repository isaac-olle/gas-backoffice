package get_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

func anIdentityProvider(t *testing.T, id ids2.IdentityProviderId, enabled bool) *identity_providers.IdentityProvider {
	t.Helper()
	idp, err := identity_providers.FromPrimitives(identity_providers.IdentityProviderPrimitives{
		Id:                    id.String(),
		Name:                  "google",
		Issuer:                "https://accounts.google.com",
		ClientId:              "client-123",
		AuthorizationEndpoint: "https://accounts.google.com/o/oauth2/v2/auth",
		TokenEndpoint:         "https://oauth2.googleapis.com/token",
		JwksUri:               "https://www.googleapis.com/oauth2/v3/certs",
		Scopes:                []string{"openid"},
		Enabled:               enabled,
	})
	require.NoError(t, err)
	return idp
}
