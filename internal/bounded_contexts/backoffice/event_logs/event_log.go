package event_logs

import (
	"encoding/json"
	"time"
)

// EventLog is a read projection of an event recorded by the audit subscriber.
type EventLog struct {
	EventId       string          `json:"event_id"`
	EventType     string          `json:"event_type"`
	AggregateId   string          `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	AuthorId      *string         `json:"author_id"`
	AuthorRole    *string         `json:"author_role"`
	Payload       json.RawMessage `json:"payload"`
	OccurredAt    time.Time       `json:"occurred_at"`
	RecordedAt    time.Time       `json:"recorded_at"`
}
