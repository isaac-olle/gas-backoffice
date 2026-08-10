package session

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	protocolEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/protocol"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
	responsetypes "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos/response_types"
)

type AuthorizationCodeSession struct {
	BaseSession
	sessionID    ids.SessionId
	redirectURI  sessionvos.RedirectURI
	state        sessionvos.State
	responseType responsetypes.ResponseType
	pkce         *sessionvos.Pkce
	sessionState sessionvos.SessionState
	consumed     sessionvos.Consumed
}

// newAuthorizationCodeSession assembles the aggregate from already-built value
// objects and records the creation event. Both the string-based public
// constructor and CreateSessionForLogin funnel through here, so a caller that
// already holds value objects never has to round-trip them back to strings.
func newAuthorizationCodeSession(
	sessionID ids.SessionId,
	clientID ids.ClientId,
	redirectURI sessionvos.RedirectURI,
	scope sessionvos.Scope,
	state sessionvos.State,
	responseType responsetypes.ResponseType,
	flow sessionvos.Flow,
	pkce *sessionvos.Pkce,
) *AuthorizationCodeSession {
	sess := &AuthorizationCodeSession{
		BaseSession:  newBaseSession(clientID, scope, flow, aggregate.AuthCodeSession),
		sessionID:    sessionID,
		redirectURI:  redirectURI,
		state:        state,
		responseType: responseType,
		pkce:         pkce,
		sessionState: sessionvos.Authorizing,
		consumed:     false,
	}
	sess.Record(protocolEvents.NewSessionCreated(
		sessionID.String(),
		clientID.String(),
		redirectURI.String(),
		scope.String(),
		flow.String(),
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
			WithCause(errorArray)
	}

	return newAuthorizationCodeSession(sessionID, clientID, redirectURI, scope, state, responseType, sessionvos.AuthorizationCodeFlow, pkce), nil
}

// CreateSessionForLogin builds an AuthorizationCodeSession for a federated
// login. The client id, redirect uri and scope are injected (they identify the
// backoffice frontend as a gas OAuth client); the response type is fixed to
// "code", the state is generated here as a value object, and there is no PKCE.
func CreateSessionForLogin(
	sessionID ids.SessionId,
	backofficeClientID string,
	redirectURI string,
	scope string,
	codeGenerator responsetypes.CodeGenerator,
) (*AuthorizationCodeSession, error) {
	const responseTypeCode = "code"
	errorArray := app_errors.NewEmptyErrorArray()

	clientID, err := ids.NewClientIdFromString(backofficeClientID)
	if err != nil {
		errorArray.AddError(err)
	}

	redirectURIVO, err := sessionvos.NewRedirectURI(redirectURI)
	if err != nil {
		errorArray.AddError(err)
	}

	scopeVO, err := sessionvos.NewScope(scope)
	if err != nil {
		errorArray.AddError(err)
	}

	state, err := sessionvos.GenerateState()
	if err != nil {
		errorArray.AddError(err)
	}

	responseType, err := responsetypes.New(responseTypeCode, codeGenerator)
	if err != nil {
		errorArray.AddError(err)
	}

	if errorArray.HasErrors() {
		return nil, app_errors.BadRequestError(ErrInvalidSessionParameters).
			WithDescription("several errors building login session").
			WithCause(errorArray)
	}

	return newAuthorizationCodeSession(sessionID, clientID, redirectURIVO, scopeVO, state, responseType, sessionvos.AuthorizationCodeFlow, nil), nil
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

	flow, err := sessionvos.NewFlowFromString(p.Flow)
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

	if errorArray.HasErrors() {
		return nil, app_errors.BadRequestError(ErrInvalidSessionParameters).
			WithDescription("several errors building session from primitives").
			WithCause(errorArray)
	}

	return &AuthorizationCodeSession{
		BaseSession:  newBaseSession(clientID, scope, flow, aggregate.AuthCodeSession),
		sessionID:    sessionID,
		redirectURI:  redirectURI,
		state:        state,
		responseType: responseType,
		pkce:         pkce,
		sessionState: sessionState,
		consumed:     sessionvos.Consumed(p.Consumed),
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
		SessionID:     s.sessionID.String(),
		ClientID:      s.clientID.String(),
		RedirectURI:   s.redirectURI.String(),
		Scope:         s.scope.String(),
		State:         s.state.String(),
		ResponseTypes: responseTypes,
		Code:          code,
		Pkce:          pkcePrimitives,
		SessionState:  s.sessionState.String(),
		Flow:          s.flow.String(),
		Consumed:      bool(s.consumed),
	}
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

func (s *AuthorizationCodeSession) Consume() {
	s.consumed = true
	s.sessionState = sessionvos.SessionStateConsumed
}

func (s *AuthorizationCodeSession) BuildRedirectURI() string {
	response := responsetypes.NewAuthorizeResponse(s.redirectURI.String())
	s.responseType.AddToResponse(response)
	response.AddState(s.state.String())
	return response.Build()
}
