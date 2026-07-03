package verify_assertion

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const grantTypeJwtBearer = "urn:ietf:params:oauth:grant-type:jwt-bearer"

type Service struct {
	finder            *domain_services.AdminAssertionKeyFinder
	assertionVerifier admin_assertion_keys.AssertionVerifier
	jwtDecoder        admin_assertion_keys.JwtDecoder
	appName           string
}

func NewService(finder *domain_services.AdminAssertionKeyFinder, assertionVerifier admin_assertion_keys.AssertionVerifier, jwtDecoder admin_assertion_keys.JwtDecoder, appName string) *Service {
	return &Service{finder: finder, assertionVerifier: assertionVerifier, jwtDecoder: jwtDecoder, appName: appName}
}

func (s *Service) VerifyAssertion(ctx context.Context, qry *Query) (any, error) {
	if qry.grantType != grantTypeJwtBearer {
		return nil, app_errors.BadRequestError("UNSUPPORTED_GRANT_TYPE").
			WithDescription("unsupported grant_type: " + qry.grantType)
	}

	kid, err := s.jwtDecoder.ExtractKeyId(ctx, qry.assertion)
	if err != nil {
		return nil, app_errors.UnauthorizedError(admin_assertion_keys.ErrInvalidAdminAssertion).
			WithDescription("failed to extract key ID from assertion").
			WithCause(err)
	}

	keyId, err := ids.NewAdminAssertionKeyIdFromString(kid)
	if err != nil {
		return nil, app_errors.UnauthorizedError(admin_assertion_keys.ErrInvalidAdminAssertion).
			WithDescription("invalid key ID in assertion").
			WithCause(err)
	}

	key, err := s.finder.Find(ctx, keyId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	claims, err := s.assertionVerifier.Verify(ctx, qry.assertion, key)
	if err != nil {
		return nil, app_errors.UnauthorizedError(admin_assertion_keys.ErrInvalidAdminAssertion).
			WithDescription("assertion signature verification failed").
			WithCause(err)
	}

	tokenClaims, err := admin_assertion_keys.NewTokenClaims(claims)
	if err != nil {
		return nil, app_errors.BadRequestError(admin_assertion_keys.ErrInvalidAdminAssertion).
			WithDescription("invalid claims in assertion").
			WithCause(err)
	}

	if err = tokenClaims.CheckIssuer(key); err != nil {
		return nil, app_errors.UnauthorizedError(admin_assertion_keys.ErrInvalidAdminAssertion).
			WithDescription("invalid issuer in assertion claims").
			WithCause(err)
	}

	if err = tokenClaims.CheckAudience(s.appName); err != nil {
		return nil, app_errors.UnauthorizedError(admin_assertion_keys.ErrInvalidAdminAssertion).
			WithDescription("invalid issuer in assertion claims").
			WithCause(err)
	}

	actor, err := tokenClaims.ToActor()
	if err != nil {
		return nil, app_errors.UnauthorizedError(admin_assertion_keys.ErrInvalidAdminAssertion).
			WithDescription("failed to convert assertion claims to actor").
			WithCause(err)
	}
	return actor, nil
}
