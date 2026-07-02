package generate_jwt

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs/domain_services"
)

type Service struct {
	finder    *domain_services.KeyPairFinder
	generator key_pairs.TokenGenerator
}

func NewService(finder *domain_services.KeyPairFinder, generator key_pairs.TokenGenerator) *Service {
	return &Service{finder: finder, generator: generator}
}

func (s *Service) GenerateJWT(ctx context.Context, qry *Query) (any, error) {
	id, err := ids.NewKeyPairIdFromString(qry.keyPairId)
	if err != nil {
		return nil, app_errors.BadRequestError(key_pairs.ErrKeyPairNotFound).
			WithDescription("invalid key pair id").
			WithCause(err)
	}

	kp, err := s.finder.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	tokenString, err := s.generator.Generate(ctx, qry.claims, kp, qry.expiresIn)
	if err != nil {
		return nil, app_errors.InternalError(key_pairs.ErrKeyPairNotFound).
			WithDescription("failed to sign jwt").
			WithCause(err)
	}

	return tokenString, nil
}
