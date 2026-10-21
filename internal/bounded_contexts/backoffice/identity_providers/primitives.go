package identity_providers

type IdentityProviderPrimitives struct {
	Id                    string   `json:"id"`
	ReferenceName         string   `json:"reference_name"`
	Name                  string   `json:"name"`
	Issuer                string   `json:"issuer"`
	ClientId              string   `json:"client_id"`
	ClientSecret          string   `json:"client_secret"`
	AuthorizationEndpoint string   `json:"authorization_endpoint"`
	TokenEndpoint         string   `json:"token_endpoint"`
	UserinfoEndpoint      string   `json:"userinfo_endpoint"`
	JwksUri               string   `json:"jwks_uri"`
	Scopes                []string `json:"scopes"`
	BackofficeRedirectUri string   `json:"backoffice_redirect_uri"`
	Enabled               bool     `json:"enabled"`
}
