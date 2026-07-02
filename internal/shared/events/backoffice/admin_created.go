package backoffice

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const AdminCreatedType event.Type = "backoffice.admin.created"

type AdminCreated struct {
	Name string
	event.BaseEvent
}

func (e AdminCreated) Type() event.Type { return AdminCreatedType }

func NewAdminCreated(aggregateId, name string) AdminCreated {
	return AdminCreated{
		Name: name,
		BaseEvent: event.NewEvent(
			aggregateId,
			AdminCreatedType,
			false,
			map[string]any{"name": name},
			aggregate.Admin,
		),
	}
}
