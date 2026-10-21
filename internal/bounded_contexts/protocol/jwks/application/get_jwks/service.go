package get_jwks

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// Service builds the JWK Set the authorization server publishes, so relying
// parties can verify the tokens it signs.
type Service struct {
	signingKeyFinder *domain_services.SigningKeyFinder
}

func NewService(signingKeyFinder *domain_services.SigningKeyFinder) *Service {
	return &Service{signingKeyFinder: signingKeyFinder}
}

func (s *Service) Get(ctx context.Context) ([]map[string]any, error) {
	keys, err := s.signingKeyFinder.FindActive(ctx)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	set := make([]map[string]any, 0, len(keys))
	for _, key := range keys {
		jwk, err := key.ToJWK()
		if err != nil {
			return nil, app_errors.InternalError(jwks.ErrInvalidSigningKey).
				WithDescription("failed to render signing key as jwk").WithCause(err, false)
		}
		set = append(set, jwk)
	}
	return set, nil
}
