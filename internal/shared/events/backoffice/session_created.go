package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const SessionCreatedType event.Type = "backoffice.session.created"

type SessionCreated struct {
	UserId string
	event.BaseEvent
}

func (e SessionCreated) Type() event.Type { return SessionCreatedType }

func NewSessionCreated(aggregateId, userId string) SessionCreated {
	return SessionCreated{
		UserId: userId,
		BaseEvent: event.NewEvent(
			aggregateId,
			SessionCreatedType,
			false,
			map[string]any{"user_id": userId},
			aggregate.Session,
		),
	}
}
