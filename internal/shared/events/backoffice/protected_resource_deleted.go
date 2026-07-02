package backoffice

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const ProtectedResourceDeletedType event.Type = "backoffice.protected_resource.deleted"

type ProtectedResourceDeleted struct {
	TenantId string
	event.BaseEvent
}

func (e ProtectedResourceDeleted) Type() event.Type { return ProtectedResourceDeletedType }

func NewProtectedResourceDeleted(aggregateId, tenantId string) ProtectedResourceDeleted {
	return ProtectedResourceDeleted{
		TenantId: tenantId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ProtectedResourceDeletedType,
			false,
			map[string]any{"tenant_id": tenantId},
			aggregate.ProtectedResource,
		),
	}
}
