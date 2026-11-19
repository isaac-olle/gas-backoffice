package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ProtectedResourceURLUpdatedType event.Type = "backoffice.protected_resource.url_updated"

type ProtectedResourceURLUpdated struct {
	URL string
	event.BaseEvent
}

func (e ProtectedResourceURLUpdated) Type() event.Type { return ProtectedResourceURLUpdatedType }

func NewProtectedResourceURLUpdated(aggregateId, url string) ProtectedResourceURLUpdated {
	return ProtectedResourceURLUpdated{
		URL: url,
		BaseEvent: event.NewEvent(
			aggregateId,
			ProtectedResourceURLUpdatedType,
			false,
			map[string]any{"url": url},
			aggregate.ProtectedResource,
		),
	}
}

type protectedResourceURLUpdatedBuilder struct{}

func NewProtectedResourceURLUpdatedBuilder() event.PrimitiveEventBuilder {
	return &protectedResourceURLUpdatedBuilder{}
}

func (b *protectedResourceURLUpdatedBuilder) Build(p *event.Primitives) (event.Event, error) {
	return ProtectedResourceURLUpdated{URL: str(p.Data, "url"), BaseEvent: event.FromPrimitivesBase(p)}, nil
}
