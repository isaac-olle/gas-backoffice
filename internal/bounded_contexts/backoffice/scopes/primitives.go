package scopes

type ScopePrimitives struct {
	Id                  string
	Name                string
	ProtectedResourceId string
	TenantId            string
	AdminId             string
	AllowedClientIds    []string
	IsSystem            bool
}

type OwnershipPrimitives struct {
	ProtectedResourceId string
	TenantId            string
	AdminId             string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{
		ProtectedResourceId: o.ProtectedResourceId.String(),
		TenantId:            o.TenantId.String(),
		AdminId:             o.AdminId.String(),
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
		TenantId:            s.ownership.TenantId.String(),
		AdminId:             s.ownership.AdminId.String(),
		AllowedClientIds:    allowedClientIds,
		IsSystem:            s.isSystem,
	}
}
