package get

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
)

type Service struct {
	clientFinder *domain_services.ClientFinder
}

func NewService(clientFinder *domain_services.ClientFinder) *Service {
	return &Service{clientFinder: clientFinder}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	id, err := ids.NewClientIdFromString(qry.clientId)
	if err != nil {
		return nil, app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id").
			WithCause(err)
	}

	client, err := s.clientFinder.Find(ctx, id)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	return client, nil
}
