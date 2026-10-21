package helpers

import (
	"context"
	"fmt"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwk"
	jwksInfra "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwk/infrastructure"
	sharedJwt "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwt"
)

// OidcTokenVerifier verifies an id_token's signature and standard claims
// against whichever IdP issued it, resolving each IdP's JWKS by its
// jwks_uri.
type OidcTokenVerifier struct {
	keys jwk.KeyFetcher
}

func NewOidcTokenVerifier() *OidcTokenVerifier {
	return &OidcTokenVerifier{keys: jwksInfra.NewHttpKeyFetcher()}
}

func (v *OidcTokenVerifier) Verify(ctx context.Context, idp *protocolShared.IdentityProviderConfig, idToken string) (map[string]any, error) {
	kid, err := sharedJwt.ExtractKid(idToken)
	if err != nil {
		return nil, err
	}

	key, err := v.keys.FetchByKid(ctx, idp.JwksUri, kid)
	if err != nil {
		return nil, fmt.Errorf("resolve idp jwks key: %w", err)
	}

	claims, err := key.VerifyString(idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid id_token: %w", err)
	}

	if iss, _ := claims["iss"].(string); iss != idp.Issuer {
		return nil, fmt.Errorf("unexpected issuer: %s", iss)
	}
	if !sharedJwt.AudienceContains(claims["aud"], idp.ClientID) {
		return nil, fmt.Errorf("unexpected audience")
	}

	return claims, nil
}
