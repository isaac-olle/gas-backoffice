package helpers

import (
	"context"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	actors2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/actors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
)

func AdminIdFromActor(ctx context.Context) (ids2.AdminId, bool) {
	actor, ok := vo2.ActorFromContext[*actors2.AdminActor](ctx)
	if !ok {
		return ids2.AdminId{}, false
	}
	return actor.Id, true
}

func UserIdFromActor(ctx context.Context) (ids2.UserId, bool) {
	actor, ok := vo2.ActorFromContext[*actors2.UserActor](ctx)
	if !ok {
		return ids2.UserId{}, false
	}
	return actor.Id, true
}

func ClientIdFromActor(ctx context.Context) (ids2.ClientId, bool) {
	actor, ok := vo2.ActorFromContext[*actors2.ClientActor](ctx)
	if !ok {
		return ids2.ClientId{}, false
	}
	return actor.Id, true
}

func UserOwnershipFromContext(ctx context.Context) (users.Ownership, bool) {
	return vo2.OwnershipFromContext[users.Ownership](ctx)
}

func UserOwnershipFromActor(ctx context.Context) (users.Ownership, bool) {
	actor, ok := vo2.ActorFromContext[*actors2.UserActor](ctx)
	if !ok {
		return users.Ownership{}, false
	}
	return actor.Ownership(), true
}

func ProtectedResourceOwnershipFromContext(ctx context.Context) (protected_resources.Ownership, bool) {
	return vo2.OwnershipFromContext[protected_resources.Ownership](ctx)
}

func GenerateScopeId(_ context.Context) (ids2.ScopeId, bool) {
	return ids2.NewScopeId(), true
}
