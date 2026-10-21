package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientGrantTypeAddedType event.Type = "backoffice.client.grant_type_added"

type ClientGrantTypeAdded struct {
	GrantType string
	event.BaseEvent
}

func (e ClientGrantTypeAdded) Type() event.Type { return ClientGrantTypeAddedType }

func NewClientGrantTypeAdded(aggregateId, grantType string) ClientGrantTypeAdded {
	return ClientGrantTypeAdded{
		GrantType: grantType,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientGrantTypeAddedType,
			false,
			map[string]any{"grant_type": grantType},
			aggregate.Client,
		),
	}
}
