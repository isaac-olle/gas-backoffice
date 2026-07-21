package security

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const UserAssertionKeyCreatedType event.Type = "security.user_assertion_key.created"

type UserAssertionKeyCreated struct {
	UserId string
	event.BaseEvent
}

func (a UserAssertionKeyCreated) Type() event.Type {
	return UserAssertionKeyCreatedType
}

func NewUserAssertionKeyCreated(aggregateId, userId string) UserAssertionKeyCreated {
	return UserAssertionKeyCreated{
		UserId: userId,
		BaseEvent: event.NewEvent(
			aggregateId,
			UserAssertionKeyCreatedType,
			false,
			map[string]any{"user_id": userId},
			aggregate.UserAssertionKey,
		),
	}
}
