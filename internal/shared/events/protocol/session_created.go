package protocol

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const SessionCreatedType event.Type = "protocol.session.created"

type SessionCreated struct {
	ClientID    string
	RedirectURI string
	Scope       string
	Flow        string
	event.BaseEvent
}

func NewSessionCreated(aggregateId, clientID, redirectURI, scope, flow string) SessionCreated {
	return SessionCreated{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		Flow:        flow,
		BaseEvent: event.NewEvent(
			aggregateId,
			SessionCreatedType,
			false,
			map[string]any{
				"client_id":    clientID,
				"redirect_uri": redirectURI,
				"scope":        scope,
				"flow":         flow,
			},
			aggregate.AuthCodeSession,
		),
	}
}
