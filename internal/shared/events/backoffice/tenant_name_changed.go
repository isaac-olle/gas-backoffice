package backoffice

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const TenantNameChangedType event.Type = "backoffice.tenant.name_changed"

type TenantNameChanged struct {
	Name string
	event.BaseEvent
}

func (e TenantNameChanged) Type() event.Type { return TenantNameChangedType }

func NewTenantNameChanged(aggregateId, name string) TenantNameChanged {
	return TenantNameChanged{
		Name: name,
		BaseEvent: event.NewEvent(
			aggregateId,
			TenantNameChangedType,
			false,
			map[string]any{"name": name},
			aggregate.Tenant,
		),
	}
}
