package create_session

import (
	"context"

	protocolSharedDS "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
	responsetypes "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos/response_types"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// Service closes the authorization step once the external IDP has
// authenticated the user: it decodes the self-contained idp_state (produced
// by authorize.Service), exchanges the code with the IDP, and creates the
// AuthorizationCodeSession already active, in a single INSERT. There is no
// longer an intermediate "authorizing" row and no separate user_identity
// aggregate — the resolved identity is embedded in the session itself.
type Service struct {
	stateCodec    session.IdpStateCodec
	clientFinder  *domain_services.ClientFinder
	idpFinder     *protocolSharedDS.IdentityProviderFinder
	idpClient     session.IdentityProviderClient
	sessionRepo   session.Repository
	codeGenerator responsetypes.CodeGenerator
	eventBus      event.Bus
}

func NewService(
	stateCodec session.IdpStateCodec,
	clientFinder *domain_services.ClientFinder,
	idpFinder *protocolSharedDS.IdentityProviderFinder,
	idpClient session.IdentityProviderClient,
	sessionRepo session.Repository,
	codeGenerator responsetypes.CodeGenerator,
	eventBus event.Bus,
) *Service {
	return &Service{
		stateCodec:    stateCodec,
		clientFinder:  clientFinder,
		idpFinder:     idpFinder,
		idpClient:     idpClient,
		sessionRepo:   sessionRepo,
		codeGenerator: codeGenerator,
		eventBus:      eventBus,
	}
}

// CreateSession handles the external IDP callback: code and state come from
// the query string, idpId is the :idp_id path param (the IDP reference name).
// It returns the client's redirect_uri, carrying gas's own authorization code
// and the client's original state.
func (s *Service) CreateSession(ctx context.Context, code, state, idpId string) (string, error) {
	payload, err := s.stateCodec.Decode(state)
	if err != nil {
		return "", app_errors.BadRequestError(session.ErrInvalidIdpState).
			WithDescription("invalid or expired idp state").
			WithCause(err, false)
	}

	if payload.IdpReferenceName != idpId {
		return "", app_errors.BadRequestError(session.ErrInvalidIdpState).
			WithDescription("idp state does not match callback identity provider")
	}

	client, err := s.clientFinder.Find(ctx, payload.ClientID)
	if err != nil {
		return "", app_errors.AddTrace(err)
	}

	idp, err := s.idpFinder.FindByReferenceName(ctx, payload.IdpReferenceName)
	if err != nil {
		return "", app_errors.AddTrace(err)
	}

	authentication, err := s.idpClient.ExchangeCode(ctx, idp, code)
	if err != nil {
		return "", app_errors.InternalError(session.ErrCannotAuthenticateWithIdp).
			WithDescription("failed to exchange code with identity provider").
			WithCause(err, false)
	}

	var refreshToken *string
	if authentication.Tokens.RefreshToken != "" {
		refreshToken = &authentication.Tokens.RefreshToken
	}

	identity, err := sessionvos.NewExternalOidcUserIdentity(
		idp.ID,
		authentication.Tokens.IDToken,
		authentication.Tokens.AccessToken,
		refreshToken,
		authentication.Claims,
		authentication.Tokens.ExpiresAt,
	)
	if err != nil {
		return "", app_errors.AddTrace(err)
	}

	sessionID := ids.NewSessionId()
	sess, err := session.NewAuthorizationCodeSession(
		sessionID,
		payload.ClientID,
		payload.RedirectURI,
		payload.Scope,
		payload.ClientState,
		payload.ResponseType,
		s.codeGenerator,
		payload.CodeChallenge,
		payload.CodeChallengeMethod,
		client,
		identity,
	)
	if err != nil {
		return "", app_errors.AddTrace(err)
	}

	if err = sess.Authenticate(); err != nil {
		return "", app_errors.AddTrace(err)
	}

	if err = s.sessionRepo.Save(ctx, sess); err != nil {
		return "", app_errors.InternalError(session.ErrCannotSaveSession).
			WithDescription("repository failed to save session").
			WithCause(err, false)
	}

	s.eventBus.Publish(ctx, sess.PullDomainEvents())

	return sess.BuildRedirectURI(), nil
}
