package flows

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/token"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

const ClientCredentialsGrantType = "client_credentials"

type ClientCredentialsGrantHandler struct {
	clientFinder    *domain_services.ClientFinder
	ccSessionRepo   session.ClientCredentialsSessionRepository
	eventBus        event.Bus
	ccSessionFinder *domain_services.ClientCredentialsSessionFinder
	tokenGenerator  *domain_services.ClientCredentialsTokenResponseGenerator
	uow             persistance.UnitOfWork
}

func NewClientCredentialsGrantHandler(
	clientFinder *domain_services.ClientFinder,
	ccSessionRepo session.ClientCredentialsSessionRepository,
	eventBus event.Bus,
	ccSessionFinder *domain_services.ClientCredentialsSessionFinder,
	tokenGenerator *domain_services.ClientCredentialsTokenResponseGenerator,
	uow persistance.UnitOfWork,
) *ClientCredentialsGrantHandler {
	return &ClientCredentialsGrantHandler{
		clientFinder:    clientFinder,
		ccSessionRepo:   ccSessionRepo,
		eventBus:        eventBus,
		ccSessionFinder: ccSessionFinder,
		tokenGenerator:  tokenGenerator,
		uow:             uow,
	}
}

func (h *ClientCredentialsGrantHandler) Supports(grantType string) bool {
	return grantType == ClientCredentialsGrantType
}

func (h *ClientCredentialsGrantHandler) Execute(ctx context.Context, req token.Request) (*session.TokenResponse, error) {
	if req.IsPublic {
		return nil, app_errors.BadRequestError(session.ErrUnsupportedGrantType).
			WithDescription("unsupported grant_type: " + ClientCredentialsGrantType)
	}

	sessionID := ids.NewSessionId()
	var ccSess *session.ClientCredentialsSession

	if err := h.uow.Do(ctx, func(ctx context.Context) error {
		client, err := h.clientFinder.Find(ctx, req.ClientID)
		if err != nil {
			return err
		}

		protectedResourceId, err := client.CheckClientCredentials(req.Scope)
		if err != nil {
			return app_errors.AddTrace(err)
		}

		ccSess, err = session.NewClientCredentialsSession(sessionID, client.ClientID, client.UserId, req.Scope, protectedResourceId)
		if err != nil {
			return app_errors.AddTrace(err)
		}

		if err = h.ccSessionRepo.Save(ctx, ccSess); err != nil {
			return app_errors.InternalError(session.ErrCannotSaveSession).
				WithDescription("failed to save client credentials session").WithCause(err, false)
		}

		h.eventBus.Publish(ctx, ccSess.PullDomainEvents())
		return nil
	}); err != nil {
		return nil, app_errors.AddTrace(err)
	}

	identity := session.NewUserIdentityProjection(ccSess.ClientID().String(), "")
	return h.tokenGenerator.Generate(ctx, ccSess, identity)
}
