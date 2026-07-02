package tenants

type TenantPrimitives struct {
	Id                   string
	Name                 string
	AdminId              string
	ClientIds            []string
	ProtectedResourceIds []string
	IsSystem             bool
}

type OwnershipPrimitives struct {
	AdminId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{AdminId: o.AdminId.String()}
}

func (t *Tenant) ToPrimitives() TenantPrimitives {
	clientIds := make([]string, 0, t.clientIds.Size())
	for _, id := range t.clientIds.Values() {
		clientIds = append(clientIds, id.String())
	}
	prIds := make([]string, 0, t.protectedResourceIds.Size())
	for _, id := range t.protectedResourceIds.Values() {
		prIds = append(prIds, id.String())
	}
	return TenantPrimitives{
		Id:                   t.id.String(),
		Name:                 t.name.String(),
		AdminId:              t.ownership.AdminId.String(),
		ClientIds:            clientIds,
		ProtectedResourceIds: prIds,
		IsSystem:             t.isSystem,
	}
}
