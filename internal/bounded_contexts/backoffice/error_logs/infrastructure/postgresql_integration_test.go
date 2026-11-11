package infrastructure_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs/infrastructure"
	"github.com/isaac-olle/gas-backoffice/test/integrationtest"
)

func TestPostgresErrorLogRepository(t *testing.T) {
	t.Parallel()
	pool := integrationtest.Postgres(t)
	ctx := context.Background()
	repo := infrastructure.NewPostgresErrorLogRepository(pool)
	log := error_logs.ErrorLog{
		AuthorId: "integration-author", AuthorRole: "subscription_user",
		ErrorId: "integration-" + time.Now().UTC().Format("20060102150405.000000000"),
		Type:    "bad_request", Code: "INTEGRATION_TEST", Description: "repository round trip",
		Details: map[string]any{"component": "integration"}, Cause: []string{"root cause"},
		StackTrace: []string{"frame"}, RequestId: "integration-request", OccurredAt: ptr(time.Now().UTC()),
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM public.errors WHERE error_id=$1`, log.ErrorId) })

	require.NoError(t, repo.Save(ctx, log))
	require.NoError(t, repo.Save(ctx, log)) // Duplicate IDs are intentionally idempotent.

	latest, err := repo.Latest(ctx, 100)
	require.NoError(t, err)
	var persisted *error_logs.ErrorLog
	for i := range latest {
		if latest[i].ErrorId == log.ErrorId {
			persisted = &latest[i]
			break
		}
	}
	require.NotNil(t, persisted)
	assert.Equal(t, log.Code, persisted.Code)
	assert.Equal(t, log.RequestId, persisted.RequestId)
	assert.Equal(t, log.AuthorId, persisted.AuthorId)
	assert.Equal(t, log.AuthorRole, persisted.AuthorRole)
	assert.Equal(t, "integration", persisted.Details["component"])
}

func ptr[T any](value T) *T { return &value }
