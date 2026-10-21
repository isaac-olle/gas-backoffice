package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const AdminNameChangedType event.Type = "backoffice.admin.name_changed"

type AdminNameChanged struct {
	Name string
	event.BaseEvent
}

func (e AdminNameChanged) Type() event.Type { return AdminNameChangedType }

func NewAdminNameChanged(aggregateId, name string) AdminNameChanged {
	return AdminNameChanged{
		Name: name,
		BaseEvent: event.NewEvent(
			aggregateId,
			AdminNameChangedType,
			false,
			map[string]any{"name": name},
			aggregate.Admin,
		),
	}
}
