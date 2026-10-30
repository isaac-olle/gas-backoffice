package scopes

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

var scopeAggregateName aggregate2.Name = aggregate2.Scope

type Scope struct {
	*aggregate2.Root
	id               ids2.ScopeId
	name             vos.ScopeName
	ownership        Ownership
	allowedClientIds *vo.SetVo[ids2.ClientId]
	isSystem         bool
	enabled          bool
	version          int
	authorId         ids2.UserId
	authorRole       vo.Role
}

func (s *Scope) EnsureClientSubscription(clientSubscriptionId ids2.SubscriptionId) error {
	if s.ownership.SubscriptionId == clientSubscriptionId {
		return nil
	}
	return app_errors.ForbiddenError(ErrInvalidOwnership).
		WithDescription("client and scope do not belong to the same subscription")
}

func (s *Scope) SubscriptionId() ids2.SubscriptionId { return s.ownership.SubscriptionId }

func (s *Scope) Id() ids2.ScopeId { return s.id }

func (s *Scope) EnsureSubscriptionAccess(allowed []ids2.SubscriptionId) error {
	for _, id := range allowed {
		if id == s.ownership.SubscriptionId {
			return nil
		}
	}
	return app_errors.ForbiddenError(ErrInvalidOwnership).
		WithDescription("resource does not belong to an accessible subscription")
}
func (s *Scope) ParentProtectedResourceId() ids2.ProtectedResourceId {
	return s.ownership.ProtectedResourceId
}

func NewScope(id ids2.ScopeId, name string, ownership Ownership, isSystem bool) (*Scope, error) {
	errArr := app_errors.NewEmptyErrorArray()

	scopeName, err := vos.NewScopeName(name)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	s := &Scope{
		Root:             aggregate2.NewRoot(scopeAggregateName),
		id:               id,
		name:             scopeName,
		ownership:        ownership,
		allowedClientIds: vo.NewEmptySetVo[ids2.ClientId](),
		isSystem:         isSystem,
		enabled:          true,
		version:          1,
	}
	s.Record(backoffice.NewScopeCreated(s.id.String(), s.name.String(), s.ownership.ProtectedResourceId.String()))
	return s, nil
}

func FromPrimitives(p ScopePrimitives) (*Scope, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewScopeIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name, err := vos.NewScopeName(p.Name)
	if err != nil {
		errArr.AddError(err)
	}

	protectedResourceId, err := ids2.NewProtectedResourceIdFromString(p.ProtectedResourceId)
	if err != nil {
		errArr.AddError(err)
	}

	var ownership Ownership
	subscriptionId, err := ids2.NewSubscriptionIdFromString(p.SubscriptionId)
	if err != nil {
		errArr.AddError(err)
	} else {
		ownership = NewOwnershipFromVos(protectedResourceId, subscriptionId)
	}

	allowedClientIds := vo.NewEmptySetVo[ids2.ClientId]()
	for _, rawId := range p.AllowedClientIds {
		clientId, cErr := ids2.NewClientIdFromString(rawId)
		if cErr != nil {
			errArr.AddError(cErr)
		} else {
			allowedClientIds.Add(clientId)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}
	var authorId ids2.UserId
	if p.AuthorId != "" {
		authorId, err = ids2.NewUserIdFromString(p.AuthorId)
		if err != nil {
			errArr.AddError(err)
		}
	}
	authorRole := vo.Role("")
	if p.AuthorRole != "" {
		authorRole, err = vo.NewRole(p.AuthorRole)
		if err != nil {
			errArr.AddError(err)
		}
	}
	if errArr.HasErrors() {
		return nil, errArr
	}

	// A persisted scope always has version >= 1 (the table defaults new rows
	// to 1, and Save always writes a positive version). A zero only ever
	// comes from a hand-built ScopePrimitives{} in a test that didn't set it.
	version := p.Version
	if version == 0 {
		version = 1
	}

	s := &Scope{
		Root:             aggregate2.NewRoot(scopeAggregateName),
		id:               id,
		name:             name,
		ownership:        ownership,
		allowedClientIds: allowedClientIds,
		isSystem:         p.IsSystem,
		enabled:          p.Enabled,
		version:          version,
		authorId:         authorId,
		authorRole:       authorRole,
	}
	return s, nil
}

func (s *Scope) IsSystemResource() bool { return s.isSystem }

func (s *Scope) IsEnabled() bool { return s.enabled }

func (s *Scope) Enable() {
	s.enabled = true
	s.Record(backoffice.NewScopeEnabled(s.id.String()))
}

func (s *Scope) Disable() {
	s.enabled = false
	s.Record(backoffice.NewScopeDisabled(s.id.String()))
}

func (s *Scope) Delete() {
	s.Record(backoffice.NewScopeDeleted(s.id.String(), s.ownership.ProtectedResourceId.String()))
}

func (s *Scope) AddAllowedClient(clientId ids2.ClientId) {
	s.allowedClientIds.Add(clientId)
	s.Record(backoffice.NewScopeClientAdded(s.id.String(), clientId.String()))
}

func (s *Scope) RemoveAllowedClient(clientId ids2.ClientId) {
	s.allowedClientIds.Remove(clientId)
	s.Record(backoffice.NewScopeClientRemoved(s.id.String(), clientId.String()))
}
