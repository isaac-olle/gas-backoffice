package infrastructure_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	shared_logger "github.com/isaac-olle/gas-backoffice/internal/shared/logger"
)

type stubOutboxRepo struct {
	primitives []*event2.Primitives
	err        error
}

func (s *stubOutboxRepo) Save(_ context.Context, _ []event2.Event) error { return nil }
func (s *stubOutboxRepo) PullPending(_ context.Context, _ int) ([]*event2.Primitives, error) {
	return s.primitives, s.err
}

type capturedPublish struct {
	ctx    context.Context
	events []event2.Event
}

type spyBus struct {
	calls []capturedPublish
}

func (s *spyBus) Publish(ctx context.Context, events []event2.Event) {
	s.calls = append(s.calls, capturedPublish{ctx: ctx, events: events})
}

type stubEvent struct {
	event2.BaseEvent
}

const stubEventType event2.Type = "test.stub"

func (stubEvent) Type() event2.Type { return stubEventType }

type stubEventBuilder struct{}

func (stubEventBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return stubEvent{BaseEvent: event2.FromPrimitivesBase(p)}, nil
}

type noopLogger struct{}

func (noopLogger) Info(_ string, _ ...shared_logger.Field)  {}
func (noopLogger) Error(_ string, _ ...shared_logger.Field) {}
func (noopLogger) Fatal(_ string, _ ...shared_logger.Field) {}
func (noopLogger) Debug(_ string, _ ...shared_logger.Field) {}

func aPrimitives(id, authorId, authorRole string) *event2.Primitives {
	return &event2.Primitives{
		ID:          id,
		Type:        stubEventType.String(),
		AggregateID: "agg-1",
		AuthorId:    authorId,
		AuthorRole:  authorRole,
	}
}

func newConsumer(t *testing.T, repo event2.OutboxRepository, bus event2.Bus) *infrastructure.Consumer {
	t.Helper()
	registry := event2.NewBuilderRegistry()
	registry.Register(stubEventType, stubEventBuilder{})
	return infrastructure.NewOutboxConsumer(repo, bus, registry, 0, noopLogger{})
}

func TestOutboxConsumer_Poll_InjectsAuthorFromPrimitives(t *testing.T) {
	t.Parallel()

	repo := &stubOutboxRepo{primitives: []*event2.Primitives{
		aPrimitives("evt-1", "00000000-0000-0000-0000-000000000001", "admin"),
	}}
	bus := &spyBus{}
	consumer := newConsumer(t, repo, bus)

	consumer.Poll(context.Background())

	require.Len(t, bus.calls, 1)
	author, ok := vo.AuthorFromContext[vo.Author](bus.calls[0].ctx)
	require.True(t, ok, "expected author to be injected into published ctx")
	prims := author.ToPrimitives()
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", prims.Id)
	assert.Equal(t, "admin", prims.Role)
}

func TestOutboxConsumer_Poll_NoAuthorPassesCtxThrough(t *testing.T) {
	t.Parallel()

	repo := &stubOutboxRepo{primitives: []*event2.Primitives{
		aPrimitives("evt-1", "", ""),
	}}
	bus := &spyBus{}
	consumer := newConsumer(t, repo, bus)

	ctx := context.Background()
	consumer.Poll(ctx)

	require.Len(t, bus.calls, 1)
	_, ok := vo.AuthorFromContext[vo.Author](bus.calls[0].ctx)
	assert.False(t, ok)
}

func TestOutboxConsumer_Poll_InvalidRoleFallsBackToCtx(t *testing.T) {
	t.Parallel()

	repo := &stubOutboxRepo{primitives: []*event2.Primitives{
		aPrimitives("evt-1", "00000000-0000-0000-0000-000000000001", "not-a-role"),
	}}
	bus := &spyBus{}
	consumer := newConsumer(t, repo, bus)

	consumer.Poll(context.Background())

	require.Len(t, bus.calls, 1)
	_, ok := vo.AuthorFromContext[vo.Author](bus.calls[0].ctx)
	assert.False(t, ok)
}

func TestOutboxConsumer_Poll_PullPendingErrorSkipsPublish(t *testing.T) {
	t.Parallel()

	repo := &stubOutboxRepo{err: errors.New("db down")}
	bus := &spyBus{}
	consumer := newConsumer(t, repo, bus)

	consumer.Poll(context.Background())

	assert.Empty(t, bus.calls)
}
