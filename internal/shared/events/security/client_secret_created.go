package security

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientSecretCreatedType event.Type = "security.client_secret.created"

type ClientSecretCreated struct {
	ClientId string
	event.BaseEvent
}

func (c ClientSecretCreated) Type() event.Type {
	return ClientSecretCreatedType
}

func NewClientSecretCreated(aggregateId, clientId string) ClientSecretCreated {
	return ClientSecretCreated{
		ClientId: clientId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientSecretCreatedType,
			false,
			map[string]any{"client_id": clientId},
			aggregate.ClientSecret,
		),
	}
}
