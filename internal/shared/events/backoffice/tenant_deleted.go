package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const TenantDeletedType event.Type = "backoffice.tenant.deleted"

type TenantDeleted struct {
	AdminId string
	event.BaseEvent
}

func (e TenantDeleted) Type() event.Type { return TenantDeletedType }

func NewTenantDeleted(aggregateId, adminId string) TenantDeleted {
	return TenantDeleted{
		AdminId: adminId,
		BaseEvent: event.NewEvent(
			aggregateId,
			TenantDeletedType,
			false,
			map[string]any{"admin_id": adminId},
			aggregate.Tenant,
		),
	}
}
