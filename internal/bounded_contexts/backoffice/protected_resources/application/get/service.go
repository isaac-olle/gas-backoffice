package get

import (
	"context"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
)

type Service struct {
	protectedResourceFinder *domain_services.ProtectedResourceFinder
}

func NewService(protectedResourceFinder *domain_services.ProtectedResourceFinder) *Service {
	return &Service{protectedResourceFinder: protectedResourceFinder}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	pr, err := s.protectedResourceFinder.Find(ctx, qry.protectedResourceId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}
	return pr, nil
}
