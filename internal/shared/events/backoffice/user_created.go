package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const UserCreatedType event.Type = "backoffice.user.created"

type UserCreated struct {
	Name string
	event.BaseEvent
}

func (e UserCreated) Type() event.Type { return UserCreatedType }

func NewUserCreated(aggregateId, name string) UserCreated {
	return UserCreated{
		Name: name,
		BaseEvent: event.NewEvent(
			aggregateId,
			UserCreatedType,
			false,
			map[string]any{"name": name},
			aggregate.User,
		),
	}
}
