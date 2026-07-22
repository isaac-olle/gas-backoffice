package helpers

import (
	"context"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	authors2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
)

func UserIdFromAuthor(ctx context.Context) (ids2.UserId, bool) {
	author, ok := vo2.AuthorFromContext[*authors2.UserAuthor](ctx)
	if !ok {
		return ids2.UserId{}, false
	}
	return author.Id, true
}

func UserOwnershipFromContext(ctx context.Context) (users.Ownership, bool) {
	return vo2.OwnershipFromContext[users.Ownership](ctx)
}

func UserOwnershipFromAuthor(ctx context.Context) (users.Ownership, bool) {
	author, ok := vo2.AuthorFromContext[*authors2.UserAuthor](ctx)
	if !ok {
		return users.Ownership{}, false
	}
	return users.NewOwnership(author.Id), true
}

func ProtectedResourceOwnershipFromContext(ctx context.Context) (protected_resources.Ownership, bool) {
	return vo2.OwnershipFromContext[protected_resources.Ownership](ctx)
}

func GenerateScopeId(_ context.Context) (ids2.ScopeId, bool) {
	return ids2.NewScopeId(), true
}
