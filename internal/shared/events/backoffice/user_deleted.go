package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const UserDeletedType event.Type = "backoffice.user.deleted"

type UserDeleted struct {
	event.BaseEvent
}

func (e UserDeleted) Type() event.Type { return UserDeletedType }

func NewUserDeleted(aggregateId string) UserDeleted {
	return UserDeleted{
		BaseEvent: event.NewEvent(
			aggregateId,
			UserDeletedType,
			false,
			map[string]any{},
			aggregate.User,
		),
	}
}
