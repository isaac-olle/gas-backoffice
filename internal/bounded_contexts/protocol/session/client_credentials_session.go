package session

import (
	"time"

	app_errors "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	protocoEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/protocol"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
)

type ClientCredentialsSession struct {
	*BaseSession
	id                  ids2.SessionId
	userId              ids2.UserId
	issuedAt            time.Time
	protectedResourceId string
}

func NewClientCredentialsSession(sessionID ids2.SessionId, clientID, rawUserId, scope, protectedResourceId string) (*ClientCredentialsSession, error) {
	errorArray := app_errors.NewEmptyErrorArray()

	clientId, err := ids2.NewClientIdFromString(clientID)
	if err != nil {
		errorArray.AddError(err)
	}

	userId, err := ids2.NewUserIdFromString(rawUserId)
	if err != nil {
		errorArray.AddError(err)
	}

	scopeVO, err := sessionvos.NewScope(scope)
	if err != nil {
		errorArray.AddError(err)
	}

	if errorArray.HasErrors() {
		return nil, errorArray
	}

	sess := &ClientCredentialsSession{
		BaseSession:         newBaseSession(clientId, scopeVO, sessionvos.ClientCredentialsGrantType, aggregate2.ClientCredentialsSession, protectedResourceId),
		id:                  sessionID,
		userId:              userId,
		issuedAt:            time.Now().UTC(),
		protectedResourceId: protectedResourceId,
	}
	sess.Record(protocoEvents.NewClientCredentialsSessionCreated(
		sess.id.String(),
		clientID,
		rawUserId,
		scope,
	))
	return sess, nil
}

func (s *ClientCredentialsSession) ProtectedResourceId() string {
	return s.protectedResourceId
}

func (s *ClientCredentialsSession) TokenClaims() TokenClaims {
	return TokenClaims{
		Scope:               s.scope.String(),
		ClientID:            s.clientID.String(),
		UserID:              s.userId.String(),
		ProtectedResourceId: s.protectedResourceId,
	}
}

func ClientCredentialsSessionFromPrimitives(p *ClientCredentialsSessionPrimitives) (*ClientCredentialsSession, error) {
	id, err := ids2.NewSessionIdFromString(p.ID)
	if err != nil {
		return nil, err
	}

	clientId, err := ids2.NewClientIdFromString(p.ClientID)
	if err != nil {
		return nil, err
	}

	userId, err := ids2.NewUserIdFromString(p.UserId)
	if err != nil {
		return nil, err
	}

	scopeVO, err := sessionvos.NewScope(p.Scope)
	if err != nil {
		return nil, err
	}

	issuedAt, err := time.Parse(time.RFC3339, p.IssuedAt)
	if err != nil {
		issuedAt = time.Now().UTC()
	}

	return &ClientCredentialsSession{
		BaseSession:         newBaseSession(clientId, scopeVO, sessionvos.ClientCredentialsGrantType, aggregate2.ClientCredentialsSession, p.ProtectedResourceId),
		id:                  id,
		userId:              userId,
		issuedAt:            issuedAt,
		protectedResourceId: p.ProtectedResourceId,
	}, nil
}

func (s *ClientCredentialsSession) ToPrimitives() *ClientCredentialsSessionPrimitives {
	return &ClientCredentialsSessionPrimitives{
		BaseSessionPrimitives: s.BaseSession.ToPrimitives(),
		ID:                    s.id.String(),
		UserId:                s.userId.String(),
		IssuedAt:              s.issuedAt.Format(time.RFC3339),
	}
}
