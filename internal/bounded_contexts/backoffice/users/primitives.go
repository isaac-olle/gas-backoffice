package users

type UserPrimitives struct {
	Id                   string
	Name                 string
	Role                 string
	ClientIds            []string
	ProtectedResourceIds []string
	IsSystem             bool
}

type OwnershipPrimitives struct {
	UserId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{UserId: o.UserId.String()}
}

func (t *User) ToPrimitives() UserPrimitives {
	clientIds := make([]string, 0, t.clientIds.Size())
	for _, id := range t.clientIds.Values() {
		clientIds = append(clientIds, id.String())
	}
	prIds := make([]string, 0, t.protectedResourceIds.Size())
	for _, id := range t.protectedResourceIds.Values() {
		prIds = append(prIds, id.String())
	}
	return UserPrimitives{
		Id:                   t.id.String(),
		Name:                 t.name.String(),
		Role:                 t.role.String(),
		ClientIds:            clientIds,
		ProtectedResourceIds: prIds,
		IsSystem:             t.isSystem,
	}
}
