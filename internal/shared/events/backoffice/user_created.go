package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const UserCreatedType event.Type = "backoffice.user.created"

type UserCreated struct {
	Name    string
	AdminId string
	event.BaseEvent
}

func (e UserCreated) Type() event.Type { return UserCreatedType }

func NewUserCreated(aggregateId, name, adminId string) UserCreated {
	return UserCreated{
		Name:    name,
		AdminId: adminId,
		BaseEvent: event.NewEvent(
			aggregateId,
			UserCreatedType,
			false,
			map[string]any{"name": name, "admin_id": adminId},
			aggregate.User,
		),
	}
}
