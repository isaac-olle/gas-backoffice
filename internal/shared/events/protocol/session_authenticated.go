package protocol

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const SessionAuthenticatedType event.Type = "protocol.session.authenticated"

type SessionAuthenticated struct {
	ClientID string
	event.BaseEvent
}

func NewSessionAuthenticated(aggregateId, clientID string) SessionAuthenticated {
	return SessionAuthenticated{
		ClientID: clientID,
		BaseEvent: event.NewEvent(
			aggregateId,
			SessionAuthenticatedType,
			false,
			map[string]any{
				"client_id": clientID,
			},
			aggregate.AuthCodeSession,
		),
	}
}
