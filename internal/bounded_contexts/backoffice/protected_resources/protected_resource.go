package protected_resources

import (
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/vos"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas/internal/shared/events/backoffice"
)

var protectedResourceAggregateName aggregate2.Name = aggregate2.ProtectedResource

type ProtectedResource struct {
	*aggregate2.Root
	id        ids2.ProtectedResourceId
	name      vos.ProtectedResourceName
	ownership Ownership
	scopeIds  *vo.SetVo[ids2.ScopeId]
	isSystem  bool
}

func NewProtectedResource(prId ids2.ProtectedResourceId, name string, ownership Ownership) (*ProtectedResource, error) {
	errArr := app_errors.NewEmptyErrorArray()

	prName, err := vos.NewProtectedResourceName(name)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	pr := &ProtectedResource{
		Root:      aggregate2.NewRoot(protectedResourceAggregateName),
		id:        prId,
		name:      prName,
		ownership: ownership,
		scopeIds:  vo.NewEmptySetVo[ids2.ScopeId](),
	}
	pr.Record(backoffice.NewProtectedResourceCreated(pr.id.String(), pr.name.String(), pr.ownership.TenantId.String()))
	return pr, nil
}

func (pr *ProtectedResource) UpdateName(newName vos.ProtectedResourceName) {
	pr.name = newName
	pr.Record(backoffice.NewProtectedResourceNameUpdated(pr.id.String(), pr.name.String()))
}

func FromPrimitives(p ProtectedResourcePrimitives) (*ProtectedResource, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewProtectedResourceIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name := vos.ProtectedResourceNameFromPrimitives(p.Name)

	ownership, err := NewOwnershipFromStrings(p.TenantId, p.AdminId)
	if err != nil {
		errArr.AddError(err)
	}

	scopeIds := vo.NewEmptySetVo[ids2.ScopeId]()
	for _, rawId := range p.ScopeIds {
		scopeId, sErr := ids2.NewScopeIdFromString(rawId)
		if sErr != nil {
			errArr.AddError(sErr)
		} else {
			scopeIds.Add(scopeId)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	pr := &ProtectedResource{
		Root:      aggregate2.NewRoot(protectedResourceAggregateName),
		id:        id,
		name:      name,
		ownership: ownership,
		scopeIds:  scopeIds,
		isSystem:  p.IsSystem,
	}
	return pr, nil
}

func (pr *ProtectedResource) IsSystemResource() bool { return pr.isSystem }

func (pr *ProtectedResource) Delete() {
	pr.Record(backoffice.NewProtectedResourceDeleted(pr.id.String(), pr.ownership.TenantId.String()))
}

func (pr *ProtectedResource) CheckOwnership(o Ownership) bool {
	return pr.ownership == o
}

func (pr *ProtectedResource) AddScopeId(scopeId ids2.ScopeId) {
	pr.scopeIds.Add(scopeId)
}

func (pr *ProtectedResource) RemoveScopeId(scopeId ids2.ScopeId) {
	pr.scopeIds.Remove(scopeId)
}
