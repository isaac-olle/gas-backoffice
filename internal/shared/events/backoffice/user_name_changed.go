package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const UserNameChangedType event.Type = "backoffice.user.name_changed"

type UserNameChanged struct {
	Name string
	event.BaseEvent
}

func (e UserNameChanged) Type() event.Type { return UserNameChangedType }

func NewUserNameChanged(aggregateId, name string) UserNameChanged {
	return UserNameChanged{
		Name: name,
		BaseEvent: event.NewEvent(
			aggregateId,
			UserNameChangedType,
			false,
			map[string]any{"name": name},
			aggregate.User,
		),
	}
}
