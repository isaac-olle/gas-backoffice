package identity_providers

type IdentityProviderPrimitives struct {
	Id                    string
	Name                  string
	Issuer                string
	ClientId              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserinfoEndpoint      string
	JwksUri               string
	Scopes                []string
	Enabled               bool
}
