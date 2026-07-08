package vo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubActorVo struct {
	id   string
	role Role
}

func (s stubActorVo) EqualsRole(r Role) bool { return r == s.role }
func (s stubActorVo) ToPrimitives() *ActorPrimitives {
	return &ActorPrimitives{Id: s.id, Role: s.role.String()}
}

func TestWithActorAndActorFromContext(t *testing.T) {
	t.Parallel()
	actor := stubActorVo{id: "user-42", role: UserRole}
	ctx := context.Background()

	ctxWithActor := WithActor(ctx, actor)
	got, ok := ActorFromContext[stubActorVo](ctxWithActor)

	assert.True(t, ok)
	assert.Equal(t, "user-42", got.ToPrimitives().Id)
	assert.Equal(t, UserRole, got.role)
}

func TestActorFromContext_NotSet(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	_, ok := ActorFromContext[stubActorVo](ctx)
	assert.False(t, ok, "should return false when no actor in context")
}

func TestActorFromContext_OverwritesPrevious(t *testing.T) {
	t.Parallel()
	first := stubActorVo{id: "first", role: AdminRole}
	second := stubActorVo{id: "second", role: ClientRole}

	ctx := WithActor(context.Background(), first)
	ctx = WithActor(ctx, second)

	got, ok := ActorFromContext[stubActorVo](ctx)
	assert.True(t, ok)
	assert.Equal(t, "second", got.ToPrimitives().Id)
}
