package create_authorization_session

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	responsetypes "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos/response_types"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Service struct {
	repository    session.Repository
	clientFinder  *domain_services.ClientFinder
	eventBus      event.Bus
	codeGenerator responsetypes.CodeGenerator
}

func NewService(repository session.Repository, clientFinder *domain_services.ClientFinder, eventBus event.Bus, codeGenerator responsetypes.CodeGenerator) *Service {
	return &Service{
		repository:    repository,
		clientFinder:  clientFinder,
		eventBus:      eventBus,
		codeGenerator: codeGenerator,
	}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	client, err := s.clientFinder.Find(ctx, cmd.clientId)
	if err != nil {
		return err
	}

	if _, err = client.CheckAuthorizationRequest(cmd.responseType, cmd.redirectUri, cmd.scope); err != nil {
		return err
	}

	sess, err := session.NewAuthorizationCodeSession(
		cmd.sessionID,
		cmd.clientId,
		cmd.redirectUri,
		cmd.scope,
		cmd.state,
		cmd.responseType,
		s.codeGenerator,
		cmd.codeChallenge,
		cmd.codeChallengeMethod,
	)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.repository.Save(ctx, sess); err != nil {
		return app_errors.InternalError(session.ErrCannotSaveSession).
			WithDescription("repository failed to save session").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, sess.PullDomainEvents())
	return nil
}
