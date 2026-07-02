package admins

type AdminPrimitives struct {
	Id        string
	Name      string
	TenantIds []string
	IsSystem  bool
}

type OwnershipPrimitives struct {
	AdminId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{AdminId: o.AdminId.String()}
}

func (a *Admin) ToPrimitives() AdminPrimitives {
	tenantIds := make([]string, 0, a.tenantIds.Size())
	for _, id := range a.tenantIds.Values() {
		tenantIds = append(tenantIds, id.String())
	}
	return AdminPrimitives{
		Id:        a.id.String(),
		Name:      a.name.String(),
		TenantIds: tenantIds,
		IsSystem:  a.isSystem,
	}
}
