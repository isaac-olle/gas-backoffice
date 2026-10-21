package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/key_factory"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/golang-jwt/jwt/v5"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/domain_services"
)

type JwtAssertionVerifier struct {
	finder *domain_services.ClientAssertionKeyFinder
}

func NewJwtAssertionVerifier(finder *domain_services.ClientAssertionKeyFinder) *JwtAssertionVerifier {
	return &JwtAssertionVerifier{finder: finder}
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
	verifyingKey, err := key_factory.NewVerifyingKeyFromPEM(p.PublicKey, p.Algorithm)
	if err != nil {
		return nil, app_errors.UnauthorizedError(client_assertion_keys.ErrInvalidClientAssertion).
			WithDescription("invalid assertion key").WithCause(

			err, false)

	}

	claims, err := verifyingKey.VerifyString(assertion)
	if err != nil {
		return nil, app_errors.UnauthorizedError(client_assertion_keys.ErrInvalidClientAssertion).
			WithDescription("assertion signature verification failed").WithCause(

			fmt.Errorf("%w", err), false)

	}

	return claims, nil
}
