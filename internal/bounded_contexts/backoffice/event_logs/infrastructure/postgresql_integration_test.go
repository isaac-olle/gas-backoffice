package infrastructure_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/event_logs/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventinfra "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/infrastructure"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Verify that records written by the existing audit subscriber's repository
// can be read through the backoffice projection, including nullable authors.
func TestPostgresEventLogRepository(t *testing.T) {
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	writer := eventinfra.NewPostgresAuditRepository(pool)
	reader := infrastructure.NewPostgresEventLogRepository(pool)
	author, role := "integration-author", "system_administrator"
	older := event.AuditRecord{
		EventID: uuid.NewString(), EventType: "integration.created", AggregateID: uuid.NewString(),
		AggregateType: "integration", AuthorID: &author, AuthorRole: &role,
		Payload:    []byte(`{"component":"integration"}`),
		OccurredAt: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	newer := older
	newer.EventID = uuid.NewString()
	newer.AuthorID, newer.AuthorRole = nil, nil
	newer.OccurredAt = older.OccurredAt.Add(time.Second)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM events.audit_events WHERE event_id IN ($1, $2)`, older.EventID, newer.EventID)
	})
	require.NoError(t, writer.Save(ctx, older))
	require.NoError(t, writer.Save(ctx, newer))
	logs, err := reader.Latest(ctx, 2)
	require.NoError(t, err)
	require.Len(t, logs, 2)
	assert.Equal(t, newer.EventID, logs[0].EventId)
	assert.Nil(t, logs[0].AuthorId)
	assert.Nil(t, logs[0].AuthorRole)
	assert.Equal(t, older.EventID, logs[1].EventId)
	assert.Equal(t, older.EventType, logs[1].EventType)
	assert.Equal(t, older.AggregateID, logs[1].AggregateId)
	assert.Equal(t, older.AggregateType, logs[1].AggregateType)
	assert.Equal(t, &author, logs[1].AuthorId)
	assert.Equal(t, &role, logs[1].AuthorRole)
	assert.True(t, older.OccurredAt.Equal(logs[1].OccurredAt))
	assert.False(t, logs[1].RecordedAt.IsZero())
	assert.JSONEq(t, string(older.Payload), string(logs[1].Payload))
	body, err := json.Marshal(logs[1])
	require.NoError(t, err)
	var response map[string]any
	require.NoError(t, json.Unmarshal(body, &response))
	assert.Equal(t, map[string]any{"component": "integration"}, response["payload"])
	assert.Equal(t, older.EventID, response["event_id"])
	empty, err := reader.Latest(ctx, 0)
	require.NoError(t, err)
	assert.NotNil(t, empty)
	assert.Empty(t, empty)
}
