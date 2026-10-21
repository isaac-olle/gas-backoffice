package modules

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/protocol"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
)

func RegisterEventBuilders(registry *event.BuilderRegistry) {
	registry.Register(backoffice.AdminCreatedType, backoffice.NewAdminCreatedBuilder())
	registry.Register(backoffice.AdminNameChangedType, backoffice.NewAdminNameChangedBuilder())
	registry.Register(backoffice.UserCreatedType, backoffice.NewUserCreatedBuilder())
	registry.Register(backoffice.UserDeletedType, backoffice.NewUserDeletedBuilder())
	registry.Register(backoffice.UserNameChangedType, backoffice.NewUserNameChangedBuilder())
	registry.Register(backoffice.SessionCreatedType, backoffice.NewSessionCreatedBuilder())
	registry.Register(backoffice.SubscriptionCreatedType, backoffice.NewSubscriptionCreatedBuilder())
	registry.Register(backoffice.SubscriptionMemberAddedType, backoffice.NewSubscriptionMemberAddedBuilder())
	registry.Register(backoffice.SubscriptionMemberRemovedType, backoffice.NewSubscriptionMemberRemovedBuilder())
	registry.Register(backoffice.ClientCreatedType, backoffice.NewClientCreatedBuilder())
	registry.Register(backoffice.ClientDeletedType, backoffice.NewClientDeletedBuilder())
	registry.Register(backoffice.ClientGrantTypeAddedType, backoffice.NewClientGrantTypeAddedBuilder())
	registry.Register(backoffice.ClientRedirectUriAddedType, backoffice.NewClientRedirectUriAddedBuilder())
	registry.Register(backoffice.ClientGrantTypeRemovedType, backoffice.NewClientGrantTypeRemovedBuilder())
	registry.Register(backoffice.ProtectedResourceCreatedType, backoffice.NewProtectedResourceCreatedBuilder())
	registry.Register(backoffice.ProtectedResourceDeletedType, backoffice.NewProtectedResourceDeletedBuilder())
	registry.Register(backoffice.ProtectedResourceNameUpdatedType, backoffice.NewProtectedResourceNameUpdatedBuilder())
	registry.Register(backoffice.ScopeCreatedType, backoffice.NewScopeCreatedBuilder())
	registry.Register(backoffice.ScopeDeletedType, backoffice.NewScopeDeletedBuilder())
	registry.Register(backoffice.ScopeClientAddedType, backoffice.NewScopeClientAddedBuilder())
	registry.Register(backoffice.ScopeClientRemovedType, backoffice.NewScopeClientRemovedBuilder())
	registry.Register(backoffice.IdentityProviderCreatedType, backoffice.NewIdentityProviderCreatedBuilder())
	registry.Register(protocol.SessionCreatedType, protocol.NewSessionCreatedBuilder())
	registry.Register(protocol.SessionAuthenticatedType, protocol.NewSessionAuthenticatedBuilder())
	registry.Register(security.ClientAssertionKeyCreatedType, security.NewClientAssertionKeyCreatedBuilder())
	registry.Register(security.ClientSecretCreatedType, security.NewClientSecretCreatedBuilder())
	registry.Register(security.KeyPairCreatedType, security.NewKeyPairCreatedBuilder())
}
