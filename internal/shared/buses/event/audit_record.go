package event

import "time"

type AuditRecord struct {
	EventID       string
	EventType     string
	AggregateID   string
	AggregateType string
	ActorID       *string
	ActorRole     *string
	Payload       []byte
	OccurredAt    time.Time
}
