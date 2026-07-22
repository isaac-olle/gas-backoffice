package vo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubAuthorVo struct {
	id   string
	role Role
}

func (s stubAuthorVo) EqualsRole(r Role) bool { return r == s.role }
func (s stubAuthorVo) ToPrimitives() *AuthorPrimitives {
	return &AuthorPrimitives{Id: s.id, Role: s.role.String()}
}

func TestWithAuthorAndAuthorFromContext(t *testing.T) {
	t.Parallel()
	author := stubAuthorVo{id: "user-42", role: TenantRole}
	ctx := context.Background()

	ctxWithAuthor := WithAuthor(ctx, author)
	got, ok := AuthorFromContext[stubAuthorVo](ctxWithAuthor)

	assert.True(t, ok)
	assert.Equal(t, "user-42", got.ToPrimitives().Id)
	assert.Equal(t, TenantRole, got.role)
}

func TestAuthorFromContext_NotSet(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	_, ok := AuthorFromContext[stubAuthorVo](ctx)
	assert.False(t, ok, "should return false when no author in context")
}

func TestAuthorFromContext_OverwritesPrevious(t *testing.T) {
	t.Parallel()
	first := stubAuthorVo{id: "first", role: AdminRole}
	second := stubAuthorVo{id: "second", role: TenantRole}

	ctx := WithAuthor(context.Background(), first)
	ctx = WithAuthor(ctx, second)

	got, ok := AuthorFromContext[stubAuthorVo](ctx)
	assert.True(t, ok)
	assert.Equal(t, "second", got.ToPrimitives().Id)
}
