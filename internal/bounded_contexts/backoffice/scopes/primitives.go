package scopes

type ScopePrimitives struct {
	Id                  string
	Name                string
	ProtectedResourceId string
	UserId              string
	AllowedClientIds    []string
	IsSystem            bool
}

type OwnershipPrimitives struct {
	ProtectedResourceId string
	UserId              string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{
		ProtectedResourceId: o.ProtectedResourceId.String(),
		UserId:              o.UserId.String(),
	}
}

func (s *Scope) ToPrimitives() ScopePrimitives {
	allowedClientIds := make([]string, 0, s.allowedClientIds.Size())
	for _, id := range s.allowedClientIds.Values() {
		allowedClientIds = append(allowedClientIds, id.String())
	}
	return ScopePrimitives{
		Id:                  s.id.String(),
		Name:                s.name.String(),
		ProtectedResourceId: s.ownership.ProtectedResourceId.String(),
		UserId:              s.ownership.UserId.String(),
		AllowedClientIds:    allowedClientIds,
		IsSystem:            s.isSystem,
	}
}
