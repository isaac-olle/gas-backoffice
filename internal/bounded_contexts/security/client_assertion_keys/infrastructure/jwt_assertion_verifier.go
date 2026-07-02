package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	jwt2 "github.com/isaac-olle/gas/internal/shared/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys/domain_services"
)

type JwtAssertionVerifier struct {
	finder   *domain_services.ClientAssertionKeyFinder
	verifier jwt2.TokenVerifier
}

func NewJwtAssertionVerifier(finder *domain_services.ClientAssertionKeyFinder, verifier jwt2.TokenVerifier) *JwtAssertionVerifier {
	return &JwtAssertionVerifier{finder: finder, verifier: verifier}
}

func (v *JwtAssertionVerifier) Verify(ctx context.Context, assertion string) (map[string]any, error) {
	unverified, _, err := jwt.NewParser().ParseUnverified(assertion, jwt.MapClaims{})
	if err != nil {
		return nil, app_errors.UnauthorizedError(client_assertion_keys.ErrInvalidClientAssertion).
			WithDescription("cannot parse assertion JWT")
	}

	raw, ok := unverified.Claims.(jwt.MapClaims)
	if !ok {
		return nil, app_errors.UnauthorizedError(client_assertion_keys.ErrInvalidClientAssertion).
			WithDescription("invalid assertion claims")
	}

	sub, ok := raw["sub"].(string)
	if !ok || sub == "" {
		return nil, app_errors.UnauthorizedError(client_assertion_keys.ErrInvalidClientAssertion).
			WithDescription("assertion must contain sub claim with client_id")
	}

	clientID, err := ids.NewClientIdFromString(sub)
	if err != nil {
		return nil, app_errors.UnauthorizedError(client_assertion_keys.ErrInvalidClientAssertion).
			WithDescription("invalid client_id in sub claim")
	}

	assertionKey, err := v.finder.FindByClientId(ctx, clientID)
	if err != nil {
		return nil, err
	}

	p := assertionKey.ToPrimitives()
	key, err := jwt2.NewKey(p.PublicKey, p.Algorithm, p.Id)
	if err != nil {
		return nil, app_errors.UnauthorizedError(client_assertion_keys.ErrInvalidClientAssertion).
			WithDescription("invalid assertion key").
			WithCause(err)
	}

	claims, err := v.verifier.Verify(assertion, key)
	if err != nil {
		return nil, app_errors.UnauthorizedError(client_assertion_keys.ErrInvalidClientAssertion).
			WithDescription("assertion signature verification failed").
			WithCause(fmt.Errorf("%w", err))
	}

	return claims, nil
}
