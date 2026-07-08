package protected_resources

type ProtectedResourcePrimitives struct {
	Id       string
	Name     string
	UserId   string
	AdminId  string
	ScopeIds []string
	IsSystem bool
}

type OwnershipPrimitives struct {
	UserId  string
	AdminId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{UserId: o.UserId.String(), AdminId: o.AdminId.String()}
}

func (pr *ProtectedResource) ToPrimitives() ProtectedResourcePrimitives {
	scopeIds := make([]string, 0, pr.scopeIds.Size())
	for _, id := range pr.scopeIds.Values() {
		scopeIds = append(scopeIds, id.String())
	}
	return ProtectedResourcePrimitives{
		Id:       pr.id.String(),
		Name:     pr.name.String(),
		UserId:   pr.ownership.UserId.String(),
		AdminId:  pr.ownership.AdminId.String(),
		ScopeIds: scopeIds,
		IsSystem: pr.isSystem,
	}
}
