package get

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type Service struct {
	adminFinder *domain_services.AdminFinder
}

func NewService(adminFinder *domain_services.AdminFinder) *Service {
	return &Service{adminFinder: adminFinder}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	admin, err := s.adminFinder.Find(ctx, qry.adminId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	if err = admin.EnsureSelf(qry.callerId); err != nil {
		return nil, app_errors.AddTrace(err)
	}

	return admin, nil
}
