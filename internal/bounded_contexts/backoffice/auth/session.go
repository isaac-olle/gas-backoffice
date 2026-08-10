package auth

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

var sessionAggregateName aggregate.Name = aggregate.Session

type Session struct {
	*aggregate.Root
	id     ids.SessionId
	userId ids.UserId
}

func NewSession(sessionId ids.SessionId, userId ids.UserId) (*Session, error) {
	s := &Session{
		Root:   aggregate.NewRoot(sessionAggregateName),
		id:     sessionId,
		userId: userId,
	}
	s.Record(backoffice.NewSessionCreated(s.id.String(), s.userId.String()))
	return s, nil
}

func FromPrimitives(p SessionPrimitives) (*Session, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids.NewSessionIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	userId, err := ids.NewUserIdFromString(p.UserId)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	return &Session{
		Root:   aggregate.NewRoot(sessionAggregateName),
		id:     id,
		userId: userId,
	}, nil
}

func (s *Session) BelongsTo(userId ids.UserId) bool {
	return s.userId == userId
}
