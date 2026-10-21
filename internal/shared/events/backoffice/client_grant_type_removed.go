package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientGrantTypeRemovedType event.Type = "backoffice.client.grant_type_removed"

type ClientGrantTypeRemoved struct {
	GrantType string
	event.BaseEvent
}

func (e ClientGrantTypeRemoved) Type() event.Type { return ClientGrantTypeRemovedType }

func NewClientGrantTypeRemoved(aggregateId, grantType string) ClientGrantTypeRemoved {
	return ClientGrantTypeRemoved{
		GrantType: grantType,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientGrantTypeRemovedType,
			false,
			map[string]any{"grant_type": grantType},
			aggregate.Client,
		),
	}
}
