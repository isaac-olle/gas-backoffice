package vo

import (
	"context"
)

type actorContextKey struct{}

/*type Actor struct {
	ID   string
	Role Role
}

func NewActor(id string, role Role) Actor {
	return Actor{ID: id, Role: role}
}

func WithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, actorContextKey{}, actor)
}

func ActorFromContext(ctx context.Context) (Actor, bool) {
	actor, ok := ctx.Value(actorContextKey{}).(Actor)
	return actor, ok
}*/

type Actor interface {
	EqualsRole(role Role) bool
	ToPrimitives() *ActorPrimitives
}

func WithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, actorContextKey{}, actor)
}

func ActorFromContext[A Actor](ctx context.Context) (A, bool) {
	actor, ok := ctx.Value(actorContextKey{}).(A)
	return actor, ok
}

type ActorPrimitives struct {
	Id   string
	Role string
}

type OwnedActor[O Ownership] interface {
	Actor
	Ownership() O
}
