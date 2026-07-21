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
	BaseSession
	id       ids2.ClientCredentialsSessionId
	userId   ids2.UserId
	issuedAt time.Time
}

func NewClientCredentialsSession(clientID, rawUserId, scope string) (*ClientCredentialsSession, error) {
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
		BaseSession: newBaseSession(clientId, scopeVO, sessionvos.NewClientCredentialsFlow(), aggregate2.ClientCredentialsSession),
		id:          ids2.NewClientCredentialsSessionId(),
		userId:      userId,
		issuedAt:    time.Now().UTC(),
	}
	sess.Record(protocoEvents.NewClientCredentialsSessionCreated(
		sess.id.String(),
		clientID,
		rawUserId,
		scope,
	))
	return sess, nil
}

func ClientCredentialsSessionFromPrimitives(p *ClientCredentialsSessionPrimitives) (*ClientCredentialsSession, error) {
	id, err := ids2.NewClientCredentialsSessionIdFromString(p.ID)
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
		BaseSession: newBaseSession(clientId, scopeVO, sessionvos.NewClientCredentialsFlow(), aggregate2.ClientCredentialsSession),
		id:          id,
		userId:      userId,
		issuedAt:    issuedAt,
	}, nil
}

func (s *ClientCredentialsSession) ToPrimitives() *ClientCredentialsSessionPrimitives {
	return &ClientCredentialsSessionPrimitives{
		ID:       s.id.String(),
		ClientID: s.clientID.String(),
		UserId:   s.userId.String(),
		Scope:    s.scope.String(),
		IssuedAt: s.issuedAt.Format(time.RFC3339),
	}
}
