package admins

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

var adminAggregateName = aggregate2.Admin

type Admin struct {
	*aggregate2.Root
	id        ids2.AdminId
	name      vos.AdminName
	ownership Ownership
	tenantIds *vo.SetVo[ids2.TenantId]
	isSystem  bool
}

func FromPrimitives(p AdminPrimitives) (*Admin, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewAdminIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name := vos.AdminNameFromPrimitives(p.Name)

	tenantIds := vo.NewEmptySetVo[ids2.TenantId]()
	for _, rawId := range p.TenantIds {
		tenantId, tErr := ids2.NewTenantIdFromString(rawId)
		if tErr != nil {
			errArr.AddError(tErr)
		} else {
			tenantIds.Add(tenantId)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	a := &Admin{
		Root:      aggregate2.NewRoot(adminAggregateName),
		id:        id,
		name:      name,
		ownership: newOwnership(id),
		tenantIds: tenantIds,
		isSystem:  p.IsSystem,
	}
	return a, nil
}

func (a *Admin) IsSystemResource() bool { return a.isSystem }

func (a *Admin) AddTenantId(tenantId ids2.TenantId) {
	a.tenantIds.Add(tenantId)
}

func (a *Admin) RemoveTenantId(tenantId ids2.TenantId) {
	a.tenantIds.Remove(tenantId)
}
