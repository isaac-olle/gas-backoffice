package security

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const KeyPairCreatedType event.Type = "security.key_pair.created"

type KeyPairCreated struct {
	Algorithm string
	event.BaseEvent
}

func (k KeyPairCreated) Type() event.Type {
	return KeyPairCreatedType
}

func NewKeyPairCreated(aggregateId, algorithm string) KeyPairCreated {
	return KeyPairCreated{
		Algorithm: algorithm,
		BaseEvent: event.NewEvent(
			aggregateId,
			KeyPairCreatedType,
			false,
			map[string]any{"algorithm": algorithm},
			aggregate.KeyPair,
		),
	}
}
