package verify_jwt

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs/domain_services"
)

type Service struct {
	finder     *domain_services.KeyPairFinder
	verifier   key_pairs.TokenVerifier
	jwtDecoder key_pairs.JwtDecoder
	appName    string
}

func NewService(finder *domain_services.KeyPairFinder, verifier key_pairs.TokenVerifier, jwtDecoder key_pairs.JwtDecoder, appName string) *Service {
	return &Service{finder: finder, verifier: verifier, jwtDecoder: jwtDecoder, appName: appName}
}

func (s *Service) VerifyJWT(ctx context.Context, qry *Query) (any, error) {
	kid, err := s.jwtDecoder.ExtractKeyId(ctx, qry.token)
	if err != nil {
		return nil, app_errors.UnauthorizedError(key_pairs.ErrInvalidToken).
			WithDescription("failed to extract key ID from token").
			WithCause(err)
	}

	keyPairId, err := ids.NewKeyPairIdFromString(kid)
	if err != nil {
		return nil, app_errors.UnauthorizedError(key_pairs.ErrInvalidToken).
			WithDescription("invalid key ID in token").
			WithCause(err)
	}

	kp, err := s.finder.Find(ctx, keyPairId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	claims, err := s.verifier.Verify(ctx, qry.token, kp)
	if err != nil {
		return nil, app_errors.UnauthorizedError(key_pairs.ErrInvalidToken).
			WithDescription("token verification failed").
			WithCause(err)
	}

	if err = claims.CheckIssuer(s.appName); err != nil {
		return nil, app_errors.UnauthorizedError(key_pairs.ErrInvalidToken).
			WithDescription("invalid issuer in token").
			WithCause(err)
	}

	if err = claims.CheckAudience(s.appName); err != nil {
		return nil, app_errors.UnauthorizedError(key_pairs.ErrInvalidToken).
			WithDescription("invalid audience in token").
			WithCause(err)
	}

	actor, err := claims.ToActor()
	if err != nil {
		return nil, app_errors.UnauthorizedError(key_pairs.ErrInvalidToken).
			WithDescription("invalid actor").
			WithCause(err)
	}
	return actor, nil
}
