package scopes

import (
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/vos"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas/internal/shared/events/backoffice"
)

var scopeAggregateName aggregate2.Name = aggregate2.Scope

type Scope struct {
	*aggregate2.Root
	id               ids2.ScopeId
	name             vos.ScopeName
	ownership        Ownership
	allowedClientIds *vo.SetVo[ids2.ClientId]
	isSystem         bool
}

func (s *Scope) CheckOwnership(o Ownership) bool {
	return s.ownership == o
}

func NewScope(id ids2.ScopeId, name string, ownership Ownership) (*Scope, error) {
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

	tenantId, err := ids2.NewTenantIdFromString(p.TenantId)
	if err != nil {
		errArr.AddError(err)
	}

	adminId, err := ids2.NewAdminIdFromString(p.AdminId)
	if err != nil {
		errArr.AddError(err)
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

	s := &Scope{
		Root:             aggregate2.NewRoot(scopeAggregateName),
		id:               id,
		name:             name,
		ownership:        newOwnership(protectedResourceId, tenantId, adminId),
		allowedClientIds: allowedClientIds,
		isSystem:         p.IsSystem,
	}
	return s, nil
}

func (s *Scope) IsSystemResource() bool { return s.isSystem }

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
