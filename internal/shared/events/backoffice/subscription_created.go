package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const SubscriptionCreatedType event.Type = "backoffice.subscription.created"

type SubscriptionCreated struct {
	Name, CreatedBy string
	event.BaseEvent
}

func (e SubscriptionCreated) Type() event.Type { return SubscriptionCreatedType }
func NewSubscriptionCreated(id, name, createdBy string) SubscriptionCreated {
	return SubscriptionCreated{Name: name,
		CreatedBy: createdBy,
		BaseEvent: event.NewEvent(id,
			SubscriptionCreatedType,
			false,
			map[string]any{"name": name,
				"created_by": createdBy},
			aggregate.Name("backoffice.subscription"))}
}
