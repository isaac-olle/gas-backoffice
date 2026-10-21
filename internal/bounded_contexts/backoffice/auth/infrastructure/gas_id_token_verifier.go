package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwk"
	jwksInfra "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwk/infrastructure"
	sharedJwt "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwt"
)

// GasIdTokenVerifier verifies the id_token the gas authorization server
// returns from /oidc/token. As an OIDC client of the gas, the backoffice
// verifies it the way any integrator would: resolving the signing key from
// the gas's published JWKS over HTTP, never from the federated IdP the user
// logged in with nor from the key pairs inside the system.
type GasIdTokenVerifier struct {
	keys     jwk.KeyFetcher
	jwksUri  string
	issuer   string
	clientId string
}

func NewGasIdTokenVerifier(jwksUri, issuer, clientId string) *GasIdTokenVerifier {
	return &GasIdTokenVerifier{
		keys:     jwksInfra.NewHttpKeyFetcher(),
		jwksUri:  jwksUri,
		issuer:   issuer,
		clientId: clientId,
	}
}

func (v *GasIdTokenVerifier) Verify(ctx context.Context, idToken string) (map[string]any, error) {
	kid, err := sharedJwt.ExtractKid(idToken)
	if err != nil {
		return nil, err
	}

	key, err := v.keys.FetchByKid(ctx, v.jwksUri, kid)
	if err != nil {
		return nil, fmt.Errorf("resolve gas jwks key: %w", err)
	}

	claims, err := key.VerifyString(idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid id_token: %w", err)
	}

	if iss, _ := claims["iss"].(string); iss != v.issuer {
		return nil, fmt.Errorf("unexpected issuer: %s", iss)
	}
	if !sharedJwt.AudienceContains(claims["aud"], v.clientId) {
		return nil, fmt.Errorf("unexpected audience")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing exp claim")
	}
	if time.Now().After(time.Unix(int64(exp), 0)) {
		return nil, fmt.Errorf("id_token expired")
	}

	return claims, nil
}
