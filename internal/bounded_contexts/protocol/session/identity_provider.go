package session

// IdentityProvider is a read projection of an external OIDC provider, analogous
// to the Client projection: plain public attributes the protocol BC needs to
// drive a federated login. It is populated from backoffice.oauth_identity_providers.
type IdentityProvider struct {
	ID                    string   `bson:"id"`
	Name                  string   `bson:"name"`
	Issuer                string   `bson:"issuer"`
	ClientID              string   `bson:"client_id"`
	ClientSecret          string   `bson:"client_secret"`
	AuthorizationEndpoint string   `bson:"authorization_endpoint"`
	TokenEndpoint         string   `bson:"token_endpoint"`
	UserinfoEndpoint      string   `bson:"userinfo_endpoint"`
	JwksURI               string   `bson:"jwks_uri"`
	Scopes                []string `bson:"scopes"`
}

func NewIdentityProvider(
	id, name, issuer, clientID, clientSecret,
	authorizationEndpoint, tokenEndpoint, userinfoEndpoint, jwksURI string,
	scopes []string,
) *IdentityProvider {
	return &IdentityProvider{
		ID:                    id,
		Name:                  name,
		Issuer:                issuer,
		ClientID:              clientID,
		ClientSecret:          clientSecret,
		AuthorizationEndpoint: authorizationEndpoint,
		TokenEndpoint:         tokenEndpoint,
		UserinfoEndpoint:      userinfoEndpoint,
		JwksURI:               jwksURI,
		Scopes:                scopes,
	}
}
