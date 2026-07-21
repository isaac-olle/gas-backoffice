package audit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/application/audit"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	sharedtypes "github.com/isaac-olle/gas-backoffice/internal/shared/types"
)

type stubRepo struct {
	saved *event.AuditRecord
	err   error
}

func (s *stubRepo) Save(_ context.Context, r event.AuditRecord) error {
	s.saved = &r
	return s.err
}

func newTestEvent() event.Event {
	return event.NewEvent("agg-1", "test.happened", false, map[string]any{"key": "val"}, sharedtypes.AggregateName("test"))
}

func TestSubscriber_SubscribedTo(t *testing.T) {
	t.Parallel()
	sub := audit.NewSubscriber(&stubRepo{})
	assert.Equal(t, event.AllTypes, sub.SubscribedTo())
}

func TestSubscriber_On_SavesRecord(t *testing.T) {
	t.Parallel()
	repo := &stubRepo{}
	sub := audit.NewSubscriber(repo)

	evt := newTestEvent()
	require.NoError(t, sub.On(context.Background(), evt))

	require.NotNil(t, repo.saved)
	assert.Equal(t, evt.Id(), repo.saved.EventID)
	assert.Equal(t, "test.happened", repo.saved.EventType)
	assert.Equal(t, "agg-1", repo.saved.AggregateID)
	assert.Equal(t, "test", repo.saved.AggregateType)
	assert.NotEmpty(t, repo.saved.Payload)
	assert.Nil(t, repo.saved.AuthorID)
	assert.Nil(t, repo.saved.AuthorRole)
}

func TestSubscriber_On_ExtractsAuthorFromContext(t *testing.T) {
	t.Parallel()
	repo := &stubRepo{}
	sub := audit.NewSubscriber(repo)

	author, err := authors.NewUserAuthorFromString("00000000-0000-0000-0000-000000000001", vo.AdminRole)
	require.NoError(t, err)
	ctx := vo.WithAuthor(context.Background(), author)

	require.NoError(t, sub.On(ctx, newTestEvent()))

	require.NotNil(t, repo.saved.AuthorID)
	require.NotNil(t, repo.saved.AuthorRole)
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", *repo.saved.AuthorID)
	assert.Equal(t, "admin", *repo.saved.AuthorRole)
}

func TestSubscriber_On_PropagatesRepoError(t *testing.T) {
	t.Parallel()
	repo := &stubRepo{err: errors.New("db down")}
	sub := audit.NewSubscriber(repo)

	err := sub.On(context.Background(), newTestEvent())
	assert.ErrorContains(t, err, "db down")
}
