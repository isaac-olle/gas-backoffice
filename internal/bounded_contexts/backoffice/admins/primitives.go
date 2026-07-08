package admins

type AdminPrimitives struct {
	Id       string
	Name     string
	UserIds  []string
	IsSystem bool
}

type OwnershipPrimitives struct {
	AdminId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{AdminId: o.AdminId.String()}
}

func (a *Admin) ToPrimitives() AdminPrimitives {
	userIds := make([]string, 0, a.userIds.Size())
	for _, id := range a.userIds.Values() {
		userIds = append(userIds, id.String())
	}
	return AdminPrimitives{
		Id:       a.id.String(),
		Name:     a.name.String(),
		UserIds:  userIds,
		IsSystem: a.isSystem,
	}
}
