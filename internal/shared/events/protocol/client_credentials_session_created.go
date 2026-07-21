package protocol

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientCredentialsSessionCreatedType event.Type = "protocol.client_credentials_session.created"

type ClientCredentialsSessionCreated struct {
	ClientID string
	UserId   string
	Scope    string
	event.BaseEvent
}

func NewClientCredentialsSessionCreated(aggregateId, clientID, userId, scope string) ClientCredentialsSessionCreated {
	return ClientCredentialsSessionCreated{
		ClientID: clientID,
		UserId:   userId,
		Scope:    scope,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientCredentialsSessionCreatedType,
			false,
			map[string]any{
				"client_id": clientID,
				"user_id":   userId,
				"scope":     scope,
			},
			aggregate.ClientCredentialsSession,
		),
	}
}
