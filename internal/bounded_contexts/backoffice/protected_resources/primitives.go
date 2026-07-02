package protected_resources

type ProtectedResourcePrimitives struct {
	Id       string
	Name     string
	TenantId string
	AdminId  string
	ScopeIds []string
	IsSystem bool
}

type OwnershipPrimitives struct {
	TenantId string
	AdminId  string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{TenantId: o.TenantId.String(), AdminId: o.AdminId.String()}
}

func (pr *ProtectedResource) ToPrimitives() ProtectedResourcePrimitives {
	scopeIds := make([]string, 0, pr.scopeIds.Size())
	for _, id := range pr.scopeIds.Values() {
		scopeIds = append(scopeIds, id.String())
	}
	return ProtectedResourcePrimitives{
		Id:       pr.id.String(),
		Name:     pr.name.String(),
		TenantId: pr.ownership.TenantId.String(),
		AdminId:  pr.ownership.AdminId.String(),
		ScopeIds: scopeIds,
		IsSystem: pr.isSystem,
	}
}
