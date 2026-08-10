package get

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type Service struct {
	sessionFinder *domain_services.SessionFinder
}

func NewService(sessionFinder *domain_services.SessionFinder) *Service {
	return &Service{sessionFinder: sessionFinder}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	session, err := s.sessionFinder.Find(ctx, qry.sessionId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}
	return session, nil
}
