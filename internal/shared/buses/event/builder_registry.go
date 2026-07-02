package event

import (
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
)

var ErrBuilderNotFound app_errors.ErrorCode = "EVENT_BUILDER_NOT_FOUND"

// PrimitiveEventBuilder reconstructs a typed domain event from its persisted primitives.
type PrimitiveEventBuilder interface {
	Build(p *Primitives) (Event, error)
}

// BuilderRegistry maps event types to their builders.
// Used by the outbox consumer to reconstruct typed events before dispatching to subscribers.
type BuilderRegistry struct {
	builders map[Type]PrimitiveEventBuilder
}

func NewBuilderRegistry() *BuilderRegistry {
	return &BuilderRegistry{
		builders: make(map[Type]PrimitiveEventBuilder),
	}
}

func (r *BuilderRegistry) Register(t Type, b PrimitiveEventBuilder) {
	r.builders[t] = b
}

func (r *BuilderRegistry) Build(p *Primitives) (Event, error) {
	builder, ok := r.builders[Type(p.Type)]
	if !ok {
		return nil, app_errors.InternalError(ErrBuilderNotFound).
			WithDescription(fmt.Sprintf("no builder registered for event type %s", p.Type))
	}
	return builder.Build(p)
}
