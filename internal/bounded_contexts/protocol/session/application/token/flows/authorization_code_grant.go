package flows

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/token"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

const AuthorizationCodeGrantType = "authorization_code"

type AuthorizationCodeGrantHandler struct {
	sessionFinder  *domain_services.AuthorizationCodeSessionFinder
	clientFinder   *domain_services.ClientFinder
	sessionRepo    session.Repository
	eventBus       event.Bus
	tokenGenerator *domain_services.AuthorizationCodeTokenResponseGenerator
	uow            persistance.UnitOfWork
}

func NewAuthorizationCodeGrantHandler(
	sessionFinder *domain_services.AuthorizationCodeSessionFinder,
	clientFinder *domain_services.ClientFinder,
	sessionRepo session.Repository,
	eventBus event.Bus,
	tokenGenerator *domain_services.AuthorizationCodeTokenResponseGenerator,
	uow persistance.UnitOfWork,
) *AuthorizationCodeGrantHandler {
	return &AuthorizationCodeGrantHandler{
		sessionFinder:  sessionFinder,
		clientFinder:   clientFinder,
		sessionRepo:    sessionRepo,
		eventBus:       eventBus,
		tokenGenerator: tokenGenerator,
		uow:            uow,
	}
}

func (h *AuthorizationCodeGrantHandler) Supports(grantType string) bool {
	return grantType == AuthorizationCodeGrantType
}

func (h *AuthorizationCodeGrantHandler) Execute(ctx context.Context, req token.Request) (*session.TokenResponse, error) {
	var sess *session.AuthorizationCodeSession

	if err := h.uow.Do(ctx, func(ctx context.Context) error {
		var err error
		sess, err = h.sessionFinder.FindByCode(ctx, req.Code)
		if err != nil {
			return err
		}

		client, err := h.clientFinder.Find(ctx, req.ClientID)
		if err != nil {
			return app_errors.AddTrace(err)
		}

		if err = sess.CheckExchangeParams(req.ClientID, req.RedirectURI); err != nil {
			return app_errors.AddTrace(err)
		}

		if _, err = sess.CheckPkce(req.CodeVerifier); err != nil {
			return app_errors.AddTrace(err)
		}

		if err = client.AssertPublic(req.IsPublic); err != nil {
			return app_errors.AddTrace(err)
		}

		sess.Consume()

		if err = h.sessionRepo.Save(ctx, sess); err != nil {
			return app_errors.InternalError(session.ErrCannotSaveSession).
				WithDescription("failed to consume session").WithCause(err, false)
		}

		h.eventBus.Publish(ctx, sess.PullDomainEvents())
		return nil
	}); err != nil {
		return nil, app_errors.AddTrace(err)
	}

	claims := sess.Identity().Claims()
	identity := session.NewUserIdentityProjection(claims.Sub, claims.Email, claims.Picture)

	return h.tokenGenerator.Generate(ctx, sess, identity)
}
