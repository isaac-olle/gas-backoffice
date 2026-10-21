package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const IdentityProviderCreatedType event.Type = "backoffice.identity_provider.created"

type IdentityProviderCreated struct {
	ReferenceName string
	Name          string
	event.BaseEvent
}

func NewIdentityProviderCreated(aggregateId, referenceName, name string) IdentityProviderCreated {
	return IdentityProviderCreated{
		ReferenceName: referenceName,
		Name:          name,
		BaseEvent: event.NewEvent(
			aggregateId,
			IdentityProviderCreatedType,
			false,
			map[string]any{
				"reference_name": referenceName,
				"name":           name,
			},
			aggregate.IdentityProvider,
		),
	}
}
