package auth

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"time"
)

const SessionTTL = time.Hour * 24

var sessionAggregateName aggregate.Name = aggregate.Session

type Session struct {
	sessionId       ids.SessionId
	userId          ids.UserId
	adminId         ids.UserId
	expiresAt       time.Time
	role            vo.Role
	subscriptionIds []ids.SubscriptionId
	*aggregate.Root
}

func NewSession(userIdString, roleString, adminIdString string, rawSubscriptionIdArgs ...[]string) (*Session, error) {
	sessionId := ids.NewSessionId()
	expiresAt := time.Now().Add(SessionTTL)

	errorArray := app_errors.NewEmptyErrorArray()

	userId, err := ids.NewUserIdFromString(userIdString)
	if err != nil {
		errorArray.AddError(err)
	}

	adminId, err := ids.NewUserIdFromString(adminIdString)
	if err != nil {
		errorArray.AddError(err)
	}

	role, err := vo.NewRole(roleString)
	if err != nil {
		errorArray.AddError(err)
	}
	var rawSubscriptionIds []string
	if len(rawSubscriptionIdArgs) > 0 {
		rawSubscriptionIds = rawSubscriptionIdArgs[0]
	}
	subscriptionIds := make([]ids.SubscriptionId, 0, len(rawSubscriptionIds))
	for _, raw := range rawSubscriptionIds {
		sid, e := ids.NewSubscriptionIdFromString(raw)
		if e != nil {
			errorArray.AddError(e)
		} else {
			subscriptionIds = append(subscriptionIds, sid)
		}
	}

	if errorArray.HasErrors() {
		return nil, app_errors.BadRequestError(ErrCannotCreateLoginSession).
			WithDescription("session constructor has failed").
			WithCause(errorArray, true)
	}

	s := &Session{
		Root:            aggregate.NewRoot(sessionAggregateName),
		sessionId:       sessionId,
		userId:          userId,
		adminId:         adminId,
		expiresAt:       expiresAt,
		role:            role,
		subscriptionIds: subscriptionIds,
	}
	s.Record(backoffice.NewSessionCreated(sessionId.String(), userIdString, roleString, expiresAt))
	return s, nil
}
