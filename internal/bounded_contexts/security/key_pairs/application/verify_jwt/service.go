package verify_jwt

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Service struct {
	finder     *domain_services.KeyPairFinder
	verifier   key_pairs.TokenVerifier
	jwtDecoder key_pairs.JwtDecoder
}

func NewService(
	finder *domain_services.KeyPairFinder,
	verifier key_pairs.TokenVerifier,
	jwtDecoder key_pairs.JwtDecoder,
) *Service {
	return &Service{
		finder:     finder,
		verifier:   verifier,
		jwtDecoder: jwtDecoder,
	}
}

func (s *Service) VerifyJWT(ctx context.Context, qry *Query) (any, error) {
	kid, err := s.jwtDecoder.ExtractKeyId(ctx, qry.token)
	if err != nil {
		return nil, app_errors.UnauthorizedError(key_pairs.ErrInvalidToken).
			WithDescription("failed to extract key ID from token").
			WithCause(err, false)
	}

	keyPairId, err := ids.NewKeyPairIdFromString(kid)
	if err != nil {
		return nil, app_errors.UnauthorizedError(key_pairs.ErrInvalidToken).
			WithDescription("invalid key ID in token").
			WithCause(err, false)
	}

	kp, err := s.finder.Find(ctx, keyPairId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	claims, err := s.verifier.Verify(ctx, qry.token, kp, qry.expectedClaims, qry.verifyExpiration)
	if err != nil {
		return nil, app_errors.UnauthorizedError(key_pairs.ErrInvalidToken).
			WithDescription("token verification failed").
			WithCause(err, false)
	}

	if len(qry.claimsToReturn) == 0 {
		return claims, nil
	}

	result, err := extractClaims(claims, qry.claimsToReturn)
	if err != nil {
		return nil, app_errors.UnauthorizedError(key_pairs.ErrInvalidToken).
			WithDescription(err.Error())
	}

	return result, nil
}

func extractClaims(claims map[string]any, template map[string]any) (map[string]any, error) {
	result := make(map[string]any, len(template))

	for name, sub := range template {
		value, ok := claims[name]
		if !ok {
			return nil, fmt.Errorf("missing claim: %s", name)
		}

		subTemplate, isNested := sub.(map[string]any)
		if !isNested {
			result[name] = value
			continue
		}

		nestedClaims, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("claim %s is not an object", name)
		}

		nestedResult, err := extractClaims(nestedClaims, subTemplate)
		if err != nil {
			return nil, err
		}
		result[name] = nestedResult
	}

	return result, nil
}
