package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientEnabledType event.Type = "backoffice.client.enabled"

type ClientEnabled struct {
	event.BaseEvent
}

func (e ClientEnabled) Type() event.Type { return ClientEnabledType }

func NewClientEnabled(aggregateId string) ClientEnabled {
	return ClientEnabled{
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientEnabledType,
			false,
			map[string]any{},
			aggregate.Client,
		),
	}
}
