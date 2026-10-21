package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
)

type Service struct {
	clientFinder *domain_services.ClientFinder
	repository   clients.Repository
}

func NewService(clientFinder *domain_services.ClientFinder, repository clients.Repository) *Service {
	return &Service{clientFinder: clientFinder, repository: repository}
}

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	clientId, err := ids.NewClientIdFromString(cmd.clientId)
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id").WithCause(err, false)
	}

	client, err := s.clientFinder.Find(ctx, clientId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	subscriptionId, _ := ids.SubscriptionIdFromContext(ctx)
	if err = client.EnsureOwnership([]ids.SubscriptionId{subscriptionId}); err != nil {
		return app_errors.AddTrace(err)
	}

	if client.IsSystemResource() {
		return app_errors.ForbiddenError(clients.ErrSystemResource).
			WithDescription("client is a system resource and cannot be deleted")
	}

	if err = s.repository.HardDelete(ctx, clientId); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to hard delete client").WithCause(err, false)
	}
	return nil
}
