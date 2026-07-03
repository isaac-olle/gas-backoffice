package get

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
)

type Service struct {
	scopeFinder *domain_services.ScopeFinder
}

func NewService(scopeFinder *domain_services.ScopeFinder) *Service {
	return &Service{scopeFinder: scopeFinder}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	return s.scopeFinder.Find(ctx, qry.scopeId)
}
