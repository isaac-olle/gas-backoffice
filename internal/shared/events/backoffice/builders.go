package backoffice

import (
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type subscriptionCreatedBuilder struct{}

func NewSubscriptionCreatedBuilder() event2.PrimitiveEventBuilder {
	return &subscriptionCreatedBuilder{}
}
func (b *subscriptionCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return SubscriptionCreated{Name: str(p.Data, "name"), CreatedBy: str(p.Data, "created_by"), BaseEvent: event2.FromPrimitivesBase(p)}, nil
}

type subscriptionMemberAddedBuilder struct{}

func NewSubscriptionMemberAddedBuilder() event2.PrimitiveEventBuilder {
	return &subscriptionMemberAddedBuilder{}
}
func (b *subscriptionMemberAddedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return SubscriptionMemberAdded{UserId: str(p.Data, "user_id"), BaseEvent: event2.FromPrimitivesBase(p)}, nil
}

type subscriptionMemberRemovedBuilder struct{}

func NewSubscriptionMemberRemovedBuilder() event2.PrimitiveEventBuilder {
	return &subscriptionMemberRemovedBuilder{}
}
func (b *subscriptionMemberRemovedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return SubscriptionMemberRemoved{UserId: str(p.Data, "user_id"), BaseEvent: event2.FromPrimitivesBase(p)}, nil
}

// ── Admin ───────────────────────────────────────────────────────────────────

type adminCreatedBuilder struct{}

func NewAdminCreatedBuilder() event2.PrimitiveEventBuilder { return &adminCreatedBuilder{} }
func (b *adminCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return AdminCreated{
		Name:      str(p.Data, "name"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type adminNameChangedBuilder struct{}

func NewAdminNameChangedBuilder() event2.PrimitiveEventBuilder { return &adminNameChangedBuilder{} }
func (b *adminNameChangedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return AdminNameChanged{
		Name:      str(p.Data, "name"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

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
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

// ── Client ────────────────────────────────────────────────────────────────────

type clientCreatedBuilder struct{}

func NewClientCreatedBuilder() event2.PrimitiveEventBuilder { return &clientCreatedBuilder{} }
func (b *clientCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientCreated{
		Name:           str(p.Data, "name"),
		SubscriptionId: str(p.Data, "subscription_id"),
		BaseEvent:      event2.FromPrimitivesBase(p),
	}, nil
}

type clientDeletedBuilder struct{}

func NewClientDeletedBuilder() event2.PrimitiveEventBuilder { return &clientDeletedBuilder{} }
func (b *clientDeletedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientDeleted{
		SubscriptionId: str(p.Data, "subscription_id"),
		BaseEvent:      event2.FromPrimitivesBase(p),
	}, nil
}

type clientGrantTypeAddedBuilder struct{}

func NewClientGrantTypeAddedBuilder() event2.PrimitiveEventBuilder {
	return &clientGrantTypeAddedBuilder{}
}
func (b *clientGrantTypeAddedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientGrantTypeAdded{
		GrantType: str(p.Data, "grant_type"),
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

type clientGrantTypeRemovedBuilder struct{}

func NewClientGrantTypeRemovedBuilder() event2.PrimitiveEventBuilder {
	return &clientGrantTypeRemovedBuilder{}
}
func (b *clientGrantTypeRemovedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientGrantTypeRemoved{
		GrantType: str(p.Data, "grant_type"),
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
		Name:           str(p.Data, "name"),
		SubscriptionId: str(p.Data, "subscription_id"),
		BaseEvent:      event2.FromPrimitivesBase(p),
	}, nil
}

type protectedResourceDeletedBuilder struct{}

func NewProtectedResourceDeletedBuilder() event2.PrimitiveEventBuilder {
	return &protectedResourceDeletedBuilder{}
}
func (b *protectedResourceDeletedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ProtectedResourceDeleted{
		SubscriptionId: str(p.Data, "subscription_id"),
		BaseEvent:      event2.FromPrimitivesBase(p),
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

// ── IdentityProvider ──────────────────────────────────────────────────────────

type identityProviderCreatedBuilder struct{}

func NewIdentityProviderCreatedBuilder() event2.PrimitiveEventBuilder {
	return &identityProviderCreatedBuilder{}
}
func (b *identityProviderCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return IdentityProviderCreated{
		ReferenceName: str(p.Data, "reference_name"),
		Name:          str(p.Data, "name"),
		BaseEvent:     event2.FromPrimitivesBase(p),
	}, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func str(data map[string]any, key string) string {
	v, _ := data[key].(string)
	return v
}
