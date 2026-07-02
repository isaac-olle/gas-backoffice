package security

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const ClientAssertionKeyCreatedType event.Type = "security.client_assertion_key.created"

type ClientAssertionKeyCreated struct {
	ClientId  string
	KeyPairId string
	event.BaseEvent
}

func (client ClientAssertionKeyCreated) Type() event.Type {
	return ClientAssertionKeyCreatedType
}

func NewClientAssertionKeyCreated(aggregateId, clientId string) ClientAssertionKeyCreated {
	return ClientAssertionKeyCreated{
		ClientId: clientId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientAssertionKeyCreatedType,
			false,
			map[string]any{"client_id": clientId},
			aggregate.ClientAssertionKey,
		),
	}
}
