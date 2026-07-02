package backoffice

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const ClientCreatedType event.Type = "backoffice.client.created"

type ClientCreated struct {
	Name     string
	AdminId  string
	TenantId string
	event.BaseEvent
}

func (e ClientCreated) Type() event.Type { return ClientCreatedType }

func NewClientCreated(aggregateId, name, adminId, tenantId string) ClientCreated {
	return ClientCreated{
		Name:     name,
		AdminId:  adminId,
		TenantId: tenantId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientCreatedType,
			false,
			map[string]any{"name": name, "admin_id": adminId, "tenant_id": tenantId},
			aggregate.Client,
		),
	}
}
