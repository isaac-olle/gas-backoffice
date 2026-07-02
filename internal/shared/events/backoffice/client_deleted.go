package backoffice

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const ClientDeletedType event.Type = "backoffice.client.deleted"

type ClientDeleted struct {
	TenantId string
	event.BaseEvent
}

func (e ClientDeleted) Type() event.Type { return ClientDeletedType }

func NewClientDeleted(aggregateId, tenantId string) ClientDeleted {
	return ClientDeleted{
		TenantId: tenantId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientDeletedType,
			false,
			map[string]any{"tenant_id": tenantId},
			aggregate.Client,
		),
	}
}
