package authorize

import (
	"context"
	"strings"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	protocolSharedDS "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
	responsetypes "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos/response_types"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// Request holds the raw /oidc/authorize parameters.
type Request struct {
	ResponseType        string
	ClientID            string
	RedirectURI         string
	Scope               string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
	IdpReferenceName    string
}

// Service validates an /oidc/authorize request and redirects the browser to
// the identity provider. It no longer persists anything: the authorization
// request is carried, self-contained and encrypted, in the "state" sent to
// the IDP, and the AuthorizationCodeSession is only created once the user
// comes back authenticated (see application/create_session).
type Service struct {
	clientFinder *domain_services.ClientFinder
	idpFinder    *protocolSharedDS.IdentityProviderFinder
	idpClient    session.AuthorizationURLBuilder
	stateCodec   session.IdpStateCodec
}

func NewService(
	clientFinder *domain_services.ClientFinder,
	idpFinder *protocolSharedDS.IdentityProviderFinder,
	idpClient session.AuthorizationURLBuilder,
	stateCodec session.IdpStateCodec,
) *Service {
	return &Service{
		clientFinder: clientFinder,
		idpFinder:    idpFinder,
		idpClient:    idpClient,
		stateCodec:   stateCodec,
	}
}

func (s *Service) Authorize(ctx context.Context, req Request) (string, error) {
	client, err := s.clientFinder.Find(ctx, req.ClientID)
	if err != nil {
		return "", app_errors.AddTrace(err)
	}

	errorArray := app_errors.NewEmptyErrorArray()

	redirectURI, err := sessionvos.NewRedirectURI(req.RedirectURI)
	if err != nil {
		errorArray.AddError(err)
	}

	if _, err = sessionvos.NewScope(req.Scope); err != nil {
		errorArray.AddError(err)
	}

	pkce, err := sessionvos.NewPkce(req.CodeChallenge, req.CodeChallengeMethod)
	if err != nil {
		errorArray.AddError(err)
	}

	if req.ResponseType != string(responsetypes.ResponseTypeCode) {
		errorArray.AddError(app_errors.BadRequestError(responsetypes.ErrInvalidResponseType).
			WithDescription("unsupported response_type: " + req.ResponseType))
	}

	if errorArray.HasErrors() {
		return "", app_errors.BadRequestError(session.ErrInvalidSessionParameters).
			WithDescription("several errors validating authorization request").
			WithCause(errorArray, true)
	}

	if !client.HasGrantType(sessionvos.AuthorizationCodeGrantType.String()) {
		return "", app_errors.BadRequestError(session.ErrGrantTypeNotAllowed).
			WithDescription("authorization_code grant type not allowed for this client")
	}

	if client.Public && !pkce.IsPresent() {
		return "", app_errors.BadRequestError(session.ErrPkceRequired).
			WithDescription("code_challenge is required for public clients")
	}

	if !client.HasRedirectURI(redirectURI.String()) {
		return "", app_errors.BadRequestError(session.ErrInvalidRedirectURI).
			WithDescription("redirect_uri not registered for this client").WithDetails(map[string]interface{}{
			"redirect_uri": redirectURI.String(),
		})
	}

	requestedScopes := strings.Fields(req.Scope)
	for _, scopeName := range requestedScopes {
		if !client.HasScope(scopeName) {
			return "", app_errors.BadRequestError(session.ErrInvalidScope).
				WithDescription("scope not allowed for this client: " + scopeName)
		}
	}

	if _, err = client.ResolveProtectedResourceId(requestedScopes); err != nil {
		return "", app_errors.AddTrace(err)
	}

	idp, err := s.idpFinder.FindByReferenceName(ctx, req.IdpReferenceName)
	if err != nil {
		return "", app_errors.AddTrace(err)
	}

	idpState, err := s.stateCodec.Encode(session.IdpStatePayload{
		ClientID:            req.ClientID,
		RedirectURI:         req.RedirectURI,
		Scope:               req.Scope,
		ResponseType:        req.ResponseType,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
		ClientState:         req.State,
		IdpReferenceName:    req.IdpReferenceName,
	})
	if err != nil {
		return "", app_errors.InternalError(session.ErrInvalidIdpState).
			WithDescription("failed to encode idp state").
			WithCause(err, false)
	}

	url, err := s.idpClient.BuildAuthorizationURL(idp, idpState)
	if err != nil {
		return "", app_errors.InternalError(protocolShared.ErrCannotGetIdentityProvider).
			WithDescription("failed to build identity provider authorization url").
			WithCause(err, false)
	}
	return url, nil
}
