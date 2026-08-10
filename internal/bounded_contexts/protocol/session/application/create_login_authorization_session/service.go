package create_login_authorization_session

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	responsetypes "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos/response_types"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

// Service creates and persists an AuthorizationCodeSession for a federated
// login. The login is not a separate concept: it is a plain session creation.
// The client id, redirect uri and scope are injected (they identify the
// backoffice frontend as a gas OAuth client); everything derived (response
// type, state) is computed by the CreateSessionForLogin constructor.
type Service struct {
	repository          session.Repository
	eventBus            event.Bus
	codeGenerator       responsetypes.CodeGenerator
	backofficeClientId  string
	redirectUri         string
	scope               string
	protectedResourceId string
}

func NewService(
	repository session.Repository,
	eventBus event.Bus,
	codeGenerator responsetypes.CodeGenerator,
	backofficeClientId,
	redirectUri,
	scope,
	protectedResourceId string,
) *Service {
	return &Service{
		repository:          repository,
		eventBus:            eventBus,
		codeGenerator:       codeGenerator,
		backofficeClientId:  backofficeClientId,
		redirectUri:         redirectUri,
		scope:               scope,
		protectedResourceId: protectedResourceId,
	}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	sess, err := session.CreateSessionForLogin(
		cmd.sessionID, s.backofficeClientId, s.redirectUri, s.scope, s.codeGenerator,
	)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.repository.Save(ctx, sess); err != nil {
		return app_errors.InternalError(session.ErrCannotSaveSession).
			WithDescription("repository failed to save login authorization session").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, sess.PullDomainEvents())
	return nil
}
