package shared

// IdentityProviderConfig is the protocol BC's read model of an external OIDC
// provider's configuration (not a synced/duplicated projection — it's read
// live from backoffice.identity_providers on each lookup). Each module uses
// only the fields it needs.
type IdentityProviderConfig struct {
	ID                    string
	Name                  string
	Issuer                string
	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	JwksUri               string
	Scopes                []string
	BackofficeRedirectUri string
}

func NewIdentityProviderConfig(
	id, name, issuer, clientID, clientSecret, authorizationEndpoint, tokenEndpoint, jwksUri string,
	scopes []string,
	backofficeRedirectUri string,
) *IdentityProviderConfig {
	return &IdentityProviderConfig{
		ID:                    id,
		Name:                  name,
		Issuer:                issuer,
		ClientID:              clientID,
		ClientSecret:          clientSecret,
		AuthorizationEndpoint: authorizationEndpoint,
		TokenEndpoint:         tokenEndpoint,
		JwksUri:               jwksUri,
		Scopes:                scopes,
		BackofficeRedirectUri: backofficeRedirectUri,
	}
}
