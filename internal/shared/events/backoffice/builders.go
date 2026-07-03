package backoffice

import (
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

// ── Tenant ────────────────────────────────────────────────────────────────────

type tenantCreatedBuilder struct{}

func NewTenantCreatedBuilder() event2.PrimitiveEventBuilder { return &tenantCreatedBuilder{} }
func (b *tenantCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return TenantCreated{
		Name:      str(p.Data, "name"),
		AdminId:   str(p.Data, "admin_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type tenantDeletedBuilder struct{}

func NewTenantDeletedBuilder() event2.PrimitiveEventBuilder { return &tenantDeletedBuilder{} }
func (b *tenantDeletedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return TenantDeleted{
		AdminId:   str(p.Data, "admin_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type tenantNameChangedBuilder struct{}

func NewTenantNameChangedBuilder() event2.PrimitiveEventBuilder { return &tenantNameChangedBuilder{} }
func (b *tenantNameChangedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return TenantNameChanged{
		Name:      str(p.Data, "name"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

// ── Client ────────────────────────────────────────────────────────────────────

type clientCreatedBuilder struct{}

func NewClientCreatedBuilder() event2.PrimitiveEventBuilder { return &clientCreatedBuilder{} }
func (b *clientCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientCreated{
		Name:      str(p.Data, "name"),
		AdminId:   str(p.Data, "admin_id"),
		TenantId:  str(p.Data, "tenant_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type clientDeletedBuilder struct{}

func NewClientDeletedBuilder() event2.PrimitiveEventBuilder { return &clientDeletedBuilder{} }
func (b *clientDeletedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ClientDeleted{
		TenantId:  str(p.Data, "tenant_id"),
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

// ── ProtectedResource ─────────────────────────────────────────────────────────

type protectedResourceCreatedBuilder struct{}

func NewProtectedResourceCreatedBuilder() event2.PrimitiveEventBuilder {
	return &protectedResourceCreatedBuilder{}
}
func (b *protectedResourceCreatedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ProtectedResourceCreated{
		Name:      str(p.Data, "name"),
		TenantId:  str(p.Data, "tenant_id"),
		BaseEvent: event2.FromPrimitivesBase(p),
	}, nil
}

type protectedResourceDeletedBuilder struct{}

func NewProtectedResourceDeletedBuilder() event2.PrimitiveEventBuilder {
	return &protectedResourceDeletedBuilder{}
}
func (b *protectedResourceDeletedBuilder) Build(p *event2.Primitives) (event2.Event, error) {
	return ProtectedResourceDeleted{
		TenantId:  str(p.Data, "tenant_id"),
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
