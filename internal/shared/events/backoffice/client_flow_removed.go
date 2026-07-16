package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientFlowRemovedType event.Type = "backoffice.client.flow_removed"

type ClientFlowRemoved struct {
	Flow string
	event.BaseEvent
}

func (e ClientFlowRemoved) Type() event.Type { return ClientFlowRemovedType }

func NewClientFlowRemoved(aggregateId, flow string) ClientFlowRemoved {
	return ClientFlowRemoved{
		Flow: flow,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientFlowRemovedType,
			false,
			map[string]any{"flow": flow},
			aggregate.Client,
		),
	}
}
