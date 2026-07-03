package event

import (
	sharedtypes "github.com/isaac-olle/gas-backoffice/internal/shared/types"
	"time"
)

type Primitives struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	AggregateID   string         `json:"aggregateId"`
	AggregateType string         `json:"aggregateType"`
	OccurredOn    time.Time      `json:"occurredOn"`
	External      bool           `json:"external"`
	Data          map[string]any `json:"data"`
}

// FromPrimitivesBase reconstructs a BaseEvent from primitives.
// Used by PrimitiveEventBuilders to embed into typed events.
func FromPrimitivesBase(p *Primitives) BaseEvent {
	return BaseEvent{
		id:            p.ID,
		aggregateID:   p.AggregateID,
		occurredOn:    p.OccurredOn,
		eventType:     Type(p.Type),
		data:          p.Data,
		external:      p.External,
		aggregateType: sharedtypes.AggregateName(p.AggregateType),
	}
}
