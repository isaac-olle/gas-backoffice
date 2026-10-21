package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"time"
)

const SessionCreatedType event.Type = "backoffice.session.created"

type SessionCreated struct {
	UserId    string
	Role      string
	ExpiresAt time.Time
	event.BaseEvent
}

func (e SessionCreated) Type() event.Type { return SessionCreatedType }

func NewSessionCreated(aggregateId, userId, role string, expiresAt time.Time) SessionCreated {
	return SessionCreated{
		UserId:    userId,
		Role:      role,
		ExpiresAt: expiresAt,
		BaseEvent: event.NewEvent(
			aggregateId,
			SessionCreatedType,
			false,
			map[string]any{
				"role":      role,
				"userId":    userId,
				"expiresAt": expiresAt,
			},
			aggregate.Session,
		),
	}
}
