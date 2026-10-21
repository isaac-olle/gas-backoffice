package hard_delete

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/domain_services"
)

type Service struct {
	finder     *domain_services.ClientAssertionKeyFinder
	repository client_assertion_keys.Repository
}

func NewService(finder *domain_services.ClientAssertionKeyFinder, repository client_assertion_keys.Repository) *Service {
	return &Service{finder: finder, repository: repository}
}

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	id, err := ids.NewClientAssertionKeyIdFromString(cmd.id)
	if err != nil {
		return app_errors.BadRequestError(client_assertion_keys.ErrClientAssertionKeyNotFound).
			WithDescription("invalid client assertion key id").WithCause(

			err, false)

	}

	clientId, err := ids.NewClientIdFromString(cmd.clientId)
	if err != nil {
		return app_errors.BadRequestError(client_assertion_keys.ErrInvalidOwnership).WithCause(err, false)
	}

	cak, err := s.finder.Find(ctx, id)
	if err != nil {
		var appErr *app_errors.AppError
		if errors.As(err, &appErr) {
			return err
		}
		return app_errors.InternalError(client_assertion_keys.ErrClientAssertionKeyNotFound).
			WithDescription("failed to find client assertion key").WithCause(err, false)
	}

	subscriptionId, _ := ids.SubscriptionIdFromContext(ctx)
	if err = cak.EnsureOwnership(clientId, []ids.SubscriptionId{subscriptionId}); err != nil {
		return err
	}

	if err = s.repository.HardDelete(ctx, id); err != nil {
		return app_errors.InternalError(client_assertion_keys.ErrClientAssertionKeyNotSaved).
			WithDescription("failed to hard delete client assertion key").WithCause(

			err, false)

	}
	return nil
}
