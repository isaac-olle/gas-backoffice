package modules

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas/internal/shared/events/security"
)

func RegisterEventBuilders(registry *event.BuilderRegistry) {
	registry.Register(backoffice.TenantCreatedType, backoffice.NewTenantCreatedBuilder())
	registry.Register(backoffice.TenantDeletedType, backoffice.NewTenantDeletedBuilder())
	registry.Register(backoffice.TenantNameChangedType, backoffice.NewTenantNameChangedBuilder())
	registry.Register(backoffice.ClientCreatedType, backoffice.NewClientCreatedBuilder())
	registry.Register(backoffice.ClientDeletedType, backoffice.NewClientDeletedBuilder())
	registry.Register(backoffice.ClientFlowAddedType, backoffice.NewClientFlowAddedBuilder())
	registry.Register(backoffice.ClientRedirectUriAddedType, backoffice.NewClientRedirectUriAddedBuilder())
	registry.Register(backoffice.ProtectedResourceCreatedType, backoffice.NewProtectedResourceCreatedBuilder())
	registry.Register(backoffice.ProtectedResourceDeletedType, backoffice.NewProtectedResourceDeletedBuilder())
	registry.Register(backoffice.ProtectedResourceNameUpdatedType, backoffice.NewProtectedResourceNameUpdatedBuilder())
	registry.Register(backoffice.ScopeCreatedType, backoffice.NewScopeCreatedBuilder())
	registry.Register(backoffice.ScopeDeletedType, backoffice.NewScopeDeletedBuilder())
	registry.Register(backoffice.ScopeClientAddedType, backoffice.NewScopeClientAddedBuilder())
	registry.Register(backoffice.ScopeClientRemovedType, backoffice.NewScopeClientRemovedBuilder())
	registry.Register(security.AdminAssertionKeyCreatedType, security.NewAdminAssertionKeyCreatedBuilder())
	registry.Register(security.ClientAssertionKeyCreatedType, security.NewClientAssertionKeyCreatedBuilder())
	registry.Register(security.ClientSecretCreatedType, security.NewClientSecretCreatedBuilder())
	registry.Register(security.KeyPairCreatedType, security.NewKeyPairCreatedBuilder())
}
