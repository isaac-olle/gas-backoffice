package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type SubscriptionId struct{ Uuid }

func NewSubscriptionId() SubscriptionId { return SubscriptionId{GenerateUuid()} }

func NewSubscriptionIdFromString(raw string) (SubscriptionId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return SubscriptionId{}, fmt.Errorf("invalid subscription id: %w", err)
	}
	return SubscriptionId{NewUuid(parsed)}, nil
}
