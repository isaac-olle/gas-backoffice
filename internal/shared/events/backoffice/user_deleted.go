package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const UserDeletedType event.Type = "backoffice.user.deleted"

type UserDeleted struct {
	AdminId string
	event.BaseEvent
}

func (e UserDeleted) Type() event.Type { return UserDeletedType }

func NewUserDeleted(aggregateId, adminId string) UserDeleted {
	return UserDeleted{
		AdminId: adminId,
		BaseEvent: event.NewEvent(
			aggregateId,
			UserDeletedType,
			false,
			map[string]any{"admin_id": adminId},
			aggregate.User,
		),
	}
}
