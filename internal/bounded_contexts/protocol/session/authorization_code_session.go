package session

import (
	"strings"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	protocolEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/protocol"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
	responsetypes "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos/response_types"
)

type AuthorizationCodeSession struct {
	*BaseSession
	sessionID         ids.SessionId
	redirectURI       sessionvos.RedirectURI
	state             sessionvos.State
	responseType      responsetypes.ResponseType
	pkce              *sessionvos.Pkce
	sessionState      sessionvos.SessionState
	consumed          sessionvos.Consumed
	clientOwnerUserId string
	identity          sessionvos.ExternalOidcUserIdentity
}

func newAuthorizationCodeSession(
	sessionID ids.SessionId,
	clientID ids.ClientId,
	redirectURI sessionvos.RedirectURI,
	scope sessionvos.Scope,
	state sessionvos.State,
	responseType responsetypes.ResponseType,
	grantType sessionvos.GrantType,
	pkce *sessionvos.Pkce,
	protectedResourceId string,
	clientOwnerUserId string,
	identity sessionvos.ExternalOidcUserIdentity,
) *AuthorizationCodeSession {
	sess := &AuthorizationCodeSession{
		BaseSession:       newBaseSession(clientID, scope, grantType, aggregate.AuthCodeSession, protectedResourceId),
		sessionID:         sessionID,
		redirectURI:       redirectURI,
		state:             state,
		responseType:      responseType,
		pkce:              pkce,
		sessionState:      sessionvos.Authorizing,
		consumed:          false,
		clientOwnerUserId: clientOwnerUserId,
		identity:          identity,
	}
	sess.Record(protocolEvents.NewSessionCreated(
		sessionID.String(),
		clientID.String(),
		redirectURI.String(),
		scope.String(),
		grantType.String(),
	))
	return sess
}

func NewAuthorizationCodeSession(
	sessionID ids.SessionId,
	clientIDString string,
	redirectURIString string,
	scopeString string,
	stateString string,
	responseTypeString string,
	codeGenerator responsetypes.CodeGenerator,
	pkceChallenge string,
	pkceMethod string,
	client *Client,
	identity sessionvos.ExternalOidcUserIdentity,
) (*AuthorizationCodeSession, error) {
	errorArray := app_errors.NewEmptyErrorArray()

	clientID, err := ids.NewClientIdFromString(clientIDString)
	if err != nil {
		errorArray.AddError(err)
	}

	redirectURI, err := sessionvos.NewRedirectURI(redirectURIString)
	if err != nil {
		errorArray.AddError(err)
	}

	scope, err := sessionvos.NewScope(scopeString)
	if err != nil {
		errorArray.AddError(err)
	}

	state, err := sessionvos.NewState(stateString)
	if err != nil {
		errorArray.AddError(err)
	}

	responseType, err := responsetypes.New(responseTypeString, codeGenerator)
	if err != nil {
		errorArray.AddError(err)
	}

	pkce, err := sessionvos.NewPkce(pkceChallenge, pkceMethod)
	if err != nil {
		errorArray.AddError(err)
	}

	if errorArray.HasErrors() {
		return nil, app_errors.BadRequestError(ErrInvalidSessionParameters).
			WithDescription("several errors building session").
			WithCause(errorArray, true)
	}

	if !client.HasGrantType(sessionvos.AuthorizationCodeGrantType.String()) {
		return nil, app_errors.BadRequestError(ErrGrantTypeNotAllowed).
			WithDescription("authorization_code grant type not allowed for this client")
	}

	if client.Public && !pkce.IsPresent() {
		return nil, app_errors.BadRequestError(ErrPkceRequired).
			WithDescription("code_challenge is required for public clients")
	}

	if !client.HasRedirectURI(redirectURI.String()) {
		return nil, app_errors.BadRequestError(ErrInvalidRedirectURI).
			WithDescription("redirect_uri not registered for this client").WithDetails(map[string]interface{}{
			"redirect_uri": redirectURI.String(),
		})
	}

	requestedScopes := strings.Fields(scope.String())
	for _, scopeName := range requestedScopes {
		if !client.HasScope(scopeName) {
			return nil, app_errors.BadRequestError(ErrInvalidScope).
				WithDescription("scope not allowed for this client: " + scopeName)
		}
	}

	protectedResourceId, err := client.ResolveProtectedResourceId(requestedScopes)
	if err != nil {
		return nil, err
	}

	return newAuthorizationCodeSession(sessionID, clientID, redirectURI, scope, state, responseType, sessionvos.AuthorizationCodeGrantType, pkce, protectedResourceId, client.UserId, identity), nil
}

func (s *AuthorizationCodeSession) ProtectedResourceId() string {
	return s.protectedResourceId
}

func (s *AuthorizationCodeSession) SessionID() ids.SessionId {
	return s.sessionID
}

func AuthorizationCodeSessionFromPrimitives(p *AuthorizationCodeSessionPrimitives) (*AuthorizationCodeSession, error) {
	if p == nil {
		return nil, app_errors.InternalError(ErrInvalidSessionParameters).
			WithDescription("session primitives cannot be nil")
	}

	errorArray := app_errors.NewEmptyErrorArray()

	sessionID, err := ids.NewSessionIdFromString(p.SessionID)
	if err != nil {
		errorArray.AddError(err)
	}

	clientID, err := ids.NewClientIdFromString(p.ClientID)
	if err != nil {
		errorArray.AddError(err)
	}

	redirectURI, err := sessionvos.NewRedirectURI(p.RedirectURI)
	if err != nil {
		errorArray.AddError(err)
	}

	scope, err := sessionvos.NewScope(p.Scope)
	if err != nil {
		errorArray.AddError(err)
	}

	state, err := sessionvos.NewState(p.State)
	if err != nil {
		errorArray.AddError(err)
	}

	responseType, err := responsetypes.FromPrimitives(p.ResponseTypes, p.Code)
	if err != nil {
		errorArray.AddError(err)
	}

	sessionState, err := sessionvos.NewSessionState(p.SessionState)
	if err != nil {
		errorArray.AddError(err)
	}

	grantType, err := sessionvos.NewGrantTypeFromString(p.GrantType)
	if err != nil {
		errorArray.AddError(err)
	}

	var pkceChallenge, pkceMethod string
	if p.Pkce != nil {
		pkceChallenge = p.Pkce.CodeChallenge
		pkceMethod = p.Pkce.Method
	}
	pkce, err := sessionvos.NewPkce(pkceChallenge, pkceMethod)
	if err != nil {
		errorArray.AddError(err)
	}

	identity, err := sessionvos.NewExternalOidcUserIdentity(
		p.Identity.IdpId,
		p.Identity.IDToken,
		p.Identity.AccessToken,
		p.Identity.RefreshToken,
		sessionvos.IdpClaims{Sub: p.Identity.Claims.Sub, Email: p.Identity.Claims.Email, Picture: p.Identity.Claims.Picture},
		p.Identity.ExpiresAt,
	)
	if err != nil {
		errorArray.AddError(err)
	}

	if errorArray.HasErrors() {
		return nil, app_errors.BadRequestError(ErrInvalidSessionParameters).
			WithDescription("several errors building session from primitives").
			WithCause(errorArray, false)

	}

	return &AuthorizationCodeSession{
		BaseSession:       newBaseSession(clientID, scope, grantType, aggregate.AuthCodeSession, p.ProtectedResourceId),
		sessionID:         sessionID,
		redirectURI:       redirectURI,
		state:             state,
		responseType:      responseType,
		pkce:              pkce,
		sessionState:      sessionState,
		consumed:          sessionvos.Consumed(p.Consumed),
		clientOwnerUserId: p.ClientOwnerUserId,
		identity:          identity,
	}, nil
}

func (s *AuthorizationCodeSession) ToPrimitives() *AuthorizationCodeSessionPrimitives {
	var pkcePrimitives *PkcePrimitives
	if s.pkce != nil {
		pkcePrimitives = &PkcePrimitives{
			CodeChallenge: s.pkce.Challenge(),
			Method:        s.pkce.Method().String(),
		}
	}

	code := ""
	switch rt := s.responseType.(type) {
	case *responsetypes.CodeResponse:
		code = rt.Code()
	case *responsetypes.StackResponseType:
		code = rt.Code()
	}

	var responseTypes []string
	switch rt := s.responseType.(type) {
	case *responsetypes.CodeResponse:
		responseTypes = []string{"code"}
	case *responsetypes.TokenResponse:
		responseTypes = []string{"token"}
	case *responsetypes.IDTokenResponse:
		responseTypes = []string{"id_token"}
	case *responsetypes.StackResponseType:
		types := rt.ResponseTypes()
		responseTypes = make([]string, len(types))
		for i, t := range types {
			responseTypes[i] = string(t)
		}
	}

	return &AuthorizationCodeSessionPrimitives{
		BaseSessionPrimitives: s.BaseSession.ToPrimitives(),
		SessionID:             s.sessionID.String(),
		ClientOwnerUserId:     s.clientOwnerUserId,
		RedirectURI:           s.redirectURI.String(),
		State:                 s.state.String(),
		ResponseTypes:         responseTypes,
		Code:                  code,
		Pkce:                  pkcePrimitives,
		SessionState:          s.sessionState.String(),
		Consumed:              bool(s.consumed),
		Identity: UserIdentityPrimitives{
			IdpId:        s.identity.IdpId().String(),
			IDToken:      s.identity.IDToken(),
			AccessToken:  s.identity.AccessToken(),
			RefreshToken: s.identity.RefreshToken(),
			Claims: IdpClaimsPrimitives{
				Sub:     s.identity.Claims().Sub,
				Email:   s.identity.Claims().Email,
				Picture: s.identity.Claims().Picture,
			},
			ExpiresAt: s.identity.ExpiresAt(),
		},
	}
}

// Identity returns the resolved external identity (tokens + claims) for this
// session, as obtained from the identity provider during create_session.
func (s *AuthorizationCodeSession) Identity() sessionvos.ExternalOidcUserIdentity {
	return s.identity
}

func (s *AuthorizationCodeSession) CheckExchangeParams(clientID, redirectURI string) error {
	if bool(s.consumed) {
		return app_errors.BadRequestError(ErrSessionConsumed).
			WithDescription("authorization code has already been used")
	}
	if s.clientID.String() != clientID {
		return app_errors.BadRequestError(ErrCodeNotFound).
			WithDescription("client_id does not match the authorization code")
	}
	if s.redirectURI.String() != redirectURI {
		return app_errors.BadRequestError(ErrInvalidRedirectURI).
			WithDescription("redirect_uri does not match the one used in the authorization request")
	}
	return nil
}

// CheckPkce verifies the PKCE code_verifier if the session uses PKCE.
// Returns (true, nil) on success, (true, err) on failure, (false, nil) if PKCE is not configured.
func (s *AuthorizationCodeSession) CheckPkce(codeVerifier string) (bool, error) {
	if s.pkce == nil {
		return false, nil
	}
	if codeVerifier == "" {
		return true, app_errors.BadRequestError(ErrInvalidCodeVerifier).
			WithDescription("code_verifier is required for PKCE flow")
	}
	if !s.pkce.Verify(codeVerifier) {
		return true, app_errors.BadRequestError(ErrInvalidCodeVerifier).
			WithDescription("code_verifier does not match code_challenge")
	}
	return true, nil
}

// Authenticate closes the authorization step once the external IDP has
// resolved the user: the session moves from authorizing to active.
func (s *AuthorizationCodeSession) Authenticate() error {
	if s.sessionState != sessionvos.Authorizing {
		return app_errors.BadRequestError(ErrSessionNotAuthorizing).
			WithDescription("session is not waiting for authentication").
			WithDetails(map[string]interface{}{"session_state": s.sessionState.String()})
	}
	s.sessionState = sessionvos.SessionStateActive
	s.Record(protocolEvents.NewSessionAuthenticated(s.sessionID.String(), s.clientID.String()))
	return nil
}

func (s *AuthorizationCodeSession) Consume() {
	s.consumed = true
	s.sessionState = sessionvos.SessionStateConsumed
}

func (s *AuthorizationCodeSession) IsConsumed() bool {
	return bool(s.consumed)
}

func (s *AuthorizationCodeSession) TokenClaims() TokenClaims {
	return TokenClaims{
		Scope:               s.scope.String(),
		ClientID:            s.clientID.String(),
		UserID:              s.clientOwnerUserId,
		ProtectedResourceId: s.protectedResourceId,
	}
}

func (s *AuthorizationCodeSession) BuildRedirectURI() string {
	response := responsetypes.NewAuthorizeResponse(s.redirectURI.String())
	s.responseType.AddToResponse(response)
	response.AddState(s.state.String())
	return response.Build()
}
