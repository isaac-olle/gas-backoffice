package tenants

import (
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/vos"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas/internal/shared/events/backoffice"
)

func NewTenant(tenantId ids2.TenantId, name string, ownership Ownership) (*Tenant, error) {
	errArr := app_errors.NewEmptyErrorArray()

	tenantName, err := vos.NewTenantName(name)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	t := &Tenant{
		Root:                 aggregate2.NewRoot(tenantAggregateName),
		id:                   tenantId,
		name:                 tenantName,
		ownership:            ownership,
		clientIds:            vo.NewEmptySetVo[ids2.ClientId](),
		protectedResourceIds: vo.NewEmptySetVo[ids2.ProtectedResourceId](),
	}
	t.Record(backoffice.NewTenantCreated(t.id.String(), t.name.String(), t.ownership.AdminId.String()))
	return t, nil
}

var tenantAggregateName aggregate2.Name = aggregate2.Tenant

type Tenant struct {
	*aggregate2.Root
	id                   ids2.TenantId
	name                 vos.TenantName
	ownership            Ownership
	clientIds            *vo.SetVo[ids2.ClientId]
	protectedResourceIds *vo.SetVo[ids2.ProtectedResourceId]
	isSystem             bool
}

func FromPrimitives(p TenantPrimitives) (*Tenant, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewTenantIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name := vos.TenantNameFromPrimitives(p.Name)

	adminId, err := ids2.NewAdminIdFromString(p.AdminId)
	if err != nil {
		errArr.AddError(err)
	}

	clientIds := vo.NewEmptySetVo[ids2.ClientId]()
	for _, rawId := range p.ClientIds {
		clientId, cErr := ids2.NewClientIdFromString(rawId)
		if cErr != nil {
			errArr.AddError(cErr)
		} else {
			clientIds.Add(clientId)
		}
	}

	protectedResourceIds := vo.NewEmptySetVo[ids2.ProtectedResourceId]()
	for _, rawId := range p.ProtectedResourceIds {
		prId, prErr := ids2.NewProtectedResourceIdFromString(rawId)
		if prErr != nil {
			errArr.AddError(prErr)
		} else {
			protectedResourceIds.Add(prId)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	t := &Tenant{
		Root:                 aggregate2.NewRoot(tenantAggregateName),
		id:                   id,
		name:                 name,
		ownership:            NewOwnership(adminId),
		clientIds:            clientIds,
		protectedResourceIds: protectedResourceIds,
		isSystem:             p.IsSystem,
	}
	return t, nil
}

func (t *Tenant) IsSystemResource() bool { return t.isSystem }

func (t *Tenant) CheckOwnership(o Ownership) bool {
	return t.ownership.Equals(o)
}

func (t *Tenant) ChangeName(newName vos.TenantName) error {
	t.name = newName
	t.Record(backoffice.NewTenantNameChanged(t.id.String(), t.name.String()))
	return nil
}

func (t *Tenant) AddClientId(clientId ids2.ClientId)    { t.clientIds.Add(clientId) }
func (t *Tenant) RemoveClientId(clientId ids2.ClientId) { t.clientIds.Remove(clientId) }
func (t *Tenant) AddProtectedResourceId(prId ids2.ProtectedResourceId) {
	t.protectedResourceIds.Add(prId)
}
func (t *Tenant) RemoveProtectedResourceId(prId ids2.ProtectedResourceId) {
	t.protectedResourceIds.Remove(prId)
}

func (t *Tenant) Delete() {
	t.Record(
		backoffice.NewTenantDeleted(t.id.String(), t.ownership.AdminId.String()),
	)
}
