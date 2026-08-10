package get

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type Service struct {
	finder *domain_services.IdentityProviderFinder
}

func NewService(finder *domain_services.IdentityProviderFinder) *Service {
	return &Service{finder: finder}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	idp, err := s.finder.FindEnabled(ctx, qry.identityProviderId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}
	return idp, nil
}
