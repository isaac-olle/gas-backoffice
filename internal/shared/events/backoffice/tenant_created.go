package backoffice

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const TenantCreatedType event.Type = "backoffice.tenant.created"

type TenantCreated struct {
	Name    string
	AdminId string
	event.BaseEvent
}

func (e TenantCreated) Type() event.Type { return TenantCreatedType }

func NewTenantCreated(aggregateId, name, adminId string) TenantCreated {
	return TenantCreated{
		Name:    name,
		AdminId: adminId,
		BaseEvent: event.NewEvent(
			aggregateId,
			TenantCreatedType,
			false,
			map[string]any{"name": name, "admin_id": adminId},
			aggregate.Tenant,
		),
	}
}
