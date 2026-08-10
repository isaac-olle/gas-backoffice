package backoffice

import (
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

// ── User ────────────────────────────────────────────────────────────────────

type userCreatedBuilder struct{}

func NewUserCreatedBuilder() event2.PrimitiveEventBuilder { return &userCreatedBuilder{} }
func (b *userCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return UserCreated{
		Name:      str(p.Data, "name"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type userDeletedBuilder struct{}

func NewUserDeletedBuilder() event2.PrimitiveEventBuilder { return &userDeletedBuilder{} }
func (b *userDeletedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return UserDeleted{
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type userNameChangedBuilder struct{}

func NewUserNameChangedBuilder() event2.PrimitiveEventBuilder { return &userNameChangedBuilder{} }
func (b *userNameChangedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return UserNameChanged{
		Name:      str(p.Data, "name"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

// ── Session ───────────────────────────────────────────────────────────────────

type sessionCreatedBuilder struct{}

func NewSessionCreatedBuilder() event2.PrimitiveEventBuilder { return &sessionCreatedBuilder{} }
func (b *sessionCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return SessionCreated{
		UserId:    str(p.Data, "user_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

// ── Client ────────────────────────────────────────────────────────────────────

type clientCreatedBuilder struct{}

func NewClientCreatedBuilder() event2.PrimitiveEventBuilder { return &clientCreatedBuilder{} }
func (b *clientCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientCreated{
		Name:      str(p.Data, "name"),
		UserId:    str(p.Data, "user_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type clientDeletedBuilder struct{}

func NewClientDeletedBuilder() event2.PrimitiveEventBuilder { return &clientDeletedBuilder{} }
func (b *clientDeletedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientDeleted{
		UserId:    str(p.Data, "user_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type clientFlowAddedBuilder struct{}

func NewClientFlowAddedBuilder() event2.PrimitiveEventBuilder { return &clientFlowAddedBuilder{} }
func (b *clientFlowAddedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientFlowAdded{
		Flow:      str(p.Data, "flow"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type clientRedirectUriAddedBuilder struct{}

func NewClientRedirectUriAddedBuilder() event2.PrimitiveEventBuilder {
	return &clientRedirectUriAddedBuilder{}
}
func (b *clientRedirectUriAddedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientRedirectUriAdded{
		RedirectUri: str(p.Data, "redirect_uri"),
		BaseEvent:   event2.FromPrimitivesBase(p),
	}, nil
}

type clientFlowRemovedBuilder struct{}

func NewClientFlowRemovedBuilder() event2.PrimitiveEventBuilder { return &clientFlowRemovedBuilder{} }
func (b *clientFlowRemovedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientFlowRemoved{
		Flow:      str(p.Data, "flow"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

// ── ProtectedResource ─────────────────────────────────────────────────────────

type protectedResourceCreatedBuilder struct{}

func NewProtectedResourceCreatedBuilder() event2.PrimitiveEventBuilder {
	return &protectedResourceCreatedBuilder{}
}
func (b *protectedResourceCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ProtectedResourceCreated{
		Name:      str(p.Data, "name"),
		UserId:    str(p.Data, "user_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type protectedResourceDeletedBuilder struct{}

func NewProtectedResourceDeletedBuilder() event2.PrimitiveEventBuilder {
	return &protectedResourceDeletedBuilder{}
}
func (b *protectedResourceDeletedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ProtectedResourceDeleted{
		UserId:    str(p.Data, "user_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type protectedResourceNameUpdatedBuilder struct{}

func NewProtectedResourceNameUpdatedBuilder() event2.PrimitiveEventBuilder {
	return &protectedResourceNameUpdatedBuilder{}
}
func (b *protectedResourceNameUpdatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ProtectedResourceNameUpdated{
		Name:      str(p.Data, "name"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

// ── Scope ─────────────────────────────────────────────────────────────────────

type scopeCreatedBuilder struct{}

func NewScopeCreatedBuilder() event2.PrimitiveEventBuilder { return &scopeCreatedBuilder{} }
func (b *scopeCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ScopeCreated{
		Name:                str(p.Data, "name"),
		ProtectedResourceId: str(p.Data, "protected_resource_id"),
		BaseEvent:           event2.FromPrimitivesBase(p),
	}, nil
}

type scopeDeletedBuilder struct{}

func NewScopeDeletedBuilder() event2.PrimitiveEventBuilder { return &scopeDeletedBuilder{} }
func (b *scopeDeletedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ScopeDeleted{
		ProtectedResourceId: str(p.Data, "protected_resource_id"),
		BaseEvent:           event2.FromPrimitivesBase(p),
	}, nil
}

type scopeClientAddedBuilder struct{}

func NewScopeClientAddedBuilder() event2.PrimitiveEventBuilder { return &scopeClientAddedBuilder{} }
func (b *scopeClientAddedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ScopeClientAdded{
		ClientId:  str(p.Data, "client_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type scopeClientRemovedBuilder struct{}

func NewScopeClientRemovedBuilder() event2.PrimitiveEventBuilder { return &scopeClientRemovedBuilder{} }
func (b *scopeClientRemovedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ScopeClientRemoved{
		ClientId:  str(p.Data, "client_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func str(data map[string]any, key string) string {
	v, _ := data[key].(string)
	return v
}
