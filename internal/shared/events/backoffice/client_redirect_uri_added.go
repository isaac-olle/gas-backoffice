package backoffice

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const ClientRedirectUriAddedType event.Type = "backoffice.client.redirect_uri_added"

type ClientRedirectUriAdded struct {
	RedirectUri string
	event.BaseEvent
}

func (e ClientRedirectUriAdded) Type() event.Type { return ClientRedirectUriAddedType }

func NewClientRedirectUriAdded(aggregateId, redirectUri string) ClientRedirectUriAdded {
	return ClientRedirectUriAdded{
		RedirectUri: redirectUri,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientRedirectUriAddedType,
			false,
			map[string]any{"redirect_uri": redirectUri},
			aggregate.Client,
		),
	}
}
