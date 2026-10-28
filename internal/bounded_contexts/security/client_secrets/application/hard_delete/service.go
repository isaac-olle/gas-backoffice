package hard_delete

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
)

type Service struct {
	repository client_secrets.Repository
	authorizer *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(repository client_secrets.Repository, authorizer *sharedDomainServices.SubscriptionAuthorizer) *Service {
	return &Service{repository: repository, authorizer: authorizer}
}

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	id, err := ids.NewClientSecretIdFromString(cmd.id)
	if err != nil {
		return app_errors.BadRequestError(client_secrets.ErrClientSecretNotFound).
			WithDescription("invalid client secret id").WithCause(

			err, false)

	}

	clientId, err := ids.NewClientIdFromString(cmd.clientId)
	if err != nil {
		return app_errors.BadRequestError(client_secrets.ErrInvalidOwnership).WithCause(err, false)
	}

	cs, err := s.repository.SearchById(ctx, id)
	if err != nil {
		var appErr *app_errors.AppError
		if errors.As(err, &appErr) {
			return err
		}
		return app_errors.InternalError(client_secrets.ErrClientSecretNotFound).
			WithDescription("failed to find client secret").WithCause(err, false)
	}
	if cs == nil {
		return app_errors.NotFoundError(client_secrets.ErrClientSecretNotFound).
			WithDescription("client secret not found")
	}

	if err = s.authorizer.Check(cs, cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}
	if err = cs.EnsureClientOwnership(clientId); err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.repository.HardDelete(ctx, id); err != nil {
		return app_errors.InternalError(client_secrets.ErrClientSecretNotSaved).
			WithDescription("failed to hard delete client secret").WithCause(

			err, false)

	}
	return nil
}
