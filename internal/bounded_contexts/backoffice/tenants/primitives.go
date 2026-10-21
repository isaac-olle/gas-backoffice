package tenants

type TenantPrimitives struct {
	Id              string   `json:"id"`
	AdminId         string   `json:"admin_id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	IsSystem        bool     `json:"is_system"`
	Version         int      `json:"version"`
	SubscriptionIds []string `json:"subscription_ids"`
}

type OwnershipPrimitives struct {
	UserId string `json:"user_id"`
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{UserId: o.UserId.String()}
}

func (t *Tenant) ToPrimitives() TenantPrimitives {
	ids := make([]string, 0, t.subscriptionIds.Size())
	for _, id := range t.subscriptionIds.Values() {
		ids = append(ids, id.String())
	}
	return TenantPrimitives{
		Id:              t.id.String(),
		AdminId:         t.adminId.String(),
		Name:            t.name.String(),
		Email:           t.email.String(),
		IsSystem:        t.isSystem,
		Version:         t.version,
		SubscriptionIds: ids,
	}
}
