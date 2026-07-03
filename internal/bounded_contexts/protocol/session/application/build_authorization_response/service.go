package build_authorization_response

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
)

type Service struct {
	sessionFinder *domain_services.SessionFinder
}

func NewService(sessionFinder *domain_services.SessionFinder) *Service {
	return &Service{sessionFinder: sessionFinder}
}

func (s *Service) BuildAuthorizationResponse(ctx context.Context, qry *Query) (any, error) {
	sess, err := s.sessionFinder.FindById(ctx, qry.sessionID)
	if err != nil {
		return nil, err
	}
	return sess.BuildRedirectURI(), nil
}
