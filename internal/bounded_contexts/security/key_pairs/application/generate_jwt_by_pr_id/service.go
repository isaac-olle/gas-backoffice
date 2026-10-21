package generate_jwt_by_pr_id

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Service struct {
	finder    *domain_services.KeyPairFinder
	generator key_pairs.TokenGenerator
}

func NewService(finder *domain_services.KeyPairFinder, generator key_pairs.TokenGenerator) *Service {
	return &Service{finder: finder, generator: generator}
}

func (s *Service) GenerateJWT(ctx context.Context, qry *Query) (any, error) {
	prId, err := ids.NewProtectedResourceIdFromString(qry.protectedResourceId)
	if err != nil {
		return nil, app_errors.BadRequestError(key_pairs.ErrKeyPairNotFound).
			WithDescription("invalid protected resource id").WithCause(

			err, false)

	}

	kp, err := s.finder.FindByProtectedResourceId(ctx, prId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	tokenString, err := s.generator.Generate(ctx, qry.protectedResourceId, qry.spec, kp, qry.expiresIn)
	if err != nil {
		return nil, app_errors.InternalError(key_pairs.ErrInvalidToken).
			WithDescription("failed to sign jwt").WithCause(

			err, false)

	}

	return tokenString, nil
}
