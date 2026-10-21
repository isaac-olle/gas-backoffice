package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientDisabledType event.Type = "backoffice.client.disabled"

type ClientDisabled struct {
	event.BaseEvent
}

func (e ClientDisabled) Type() event.Type { return ClientDisabledType }

func NewClientDisabled(aggregateId string) ClientDisabled {
	return ClientDisabled{
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientDisabledType,
			false,
			map[string]any{},
			aggregate.Client,
		),
	}
}
