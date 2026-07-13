package event

import "time"

type AuditRecord struct {
	EventID       string
	EventType     string
	AggregateID   string
	AggregateType string
	AuthorID      *string
	AuthorRole    *string
	Payload       []byte
	OccurredAt    time.Time
}
