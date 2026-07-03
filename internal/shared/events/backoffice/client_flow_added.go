package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientFlowAddedType event.Type = "backoffice.client.flow_added"

type ClientFlowAdded struct {
	Flow string
	event.BaseEvent
}

func (e ClientFlowAdded) Type() event.Type { return ClientFlowAddedType }

func NewClientFlowAdded(aggregateId, flow string) ClientFlowAdded {
	return ClientFlowAdded{
		Flow: flow,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientFlowAddedType,
			false,
			map[string]any{"flow": flow},
			aggregate.Client,
		),
	}
}
