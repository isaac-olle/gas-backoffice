package event

import (
	sharedtypes "github.com/isaac-olle/gas-backoffice/internal/shared/types"
	"time"

	"github.com/google/uuid"
)

type Type string

var AllTypes Type = "*"

func (t Type) String() string { return string(t) }

type Event interface {
	Id() string
	AggregateID() string
	External() bool
	Type() Type
	Data() map[string]any
	OccurredOn() time.Time
	AggregateType() sharedtypes.AggregateName
	Primitives() *Primitives
}

type BaseEvent struct {
	id            string
	aggregateID   string
	occurredOn    time.Time
	data          map[string]any
	external      bool
	eventType     Type
	aggregateType sharedtypes.AggregateName
}

func NewEvent(aggregateID string, eventType Type, external bool, data map[string]any, aggregateType sharedtypes.AggregateName) BaseEvent {
	if data == nil {
		data = make(map[string]any)
	}
	return BaseEvent{
		id:            uuid.New().String(),
		aggregateID:   aggregateID,
		occurredOn:    time.Now(),
		eventType:     eventType,
		external:      external,
		data:          data,
		aggregateType: aggregateType,
	}
}

func (b BaseEvent) Id() string                               { return b.id }
func (b BaseEvent) AggregateID() string                      { return b.aggregateID }
func (b BaseEvent) OccurredOn() time.Time                    { return b.occurredOn }
func (b BaseEvent) Data() map[string]any                     { return b.data }
func (b BaseEvent) External() bool                           { return b.external }
func (b BaseEvent) Type() Type                               { return b.eventType }
func (b BaseEvent) AggregateType() sharedtypes.AggregateName { return b.aggregateType }

func (b BaseEvent) Primitives() *Primitives {
	return &Primitives{
		ID:            b.Id(),
		Type:          b.Type().String(),
		AggregateID:   b.aggregateID,
		AggregateType: b.aggregateType.String(),
		OccurredOn:    b.OccurredOn(),
		External:      b.External(),
		Data:          b.Data(),
	}
}
