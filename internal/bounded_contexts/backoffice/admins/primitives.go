package admins

type AdminPrimitives struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	IsSystem        bool     `json:"is_system"`
	Version         int      `json:"version"`
	SubscriptionIds []string `json:"subscription_ids"`
}

func (a *Admin) ToPrimitives() AdminPrimitives {
	ids := make([]string, 0, a.subscriptionIds.Size())
	for _, id := range a.subscriptionIds.Values() {
		ids = append(ids, id.String())
	}
	return AdminPrimitives{
		Id:              a.id.String(),
		Name:            a.name.String(),
		Email:           a.email.String(),
		IsSystem:        a.isSystem,
		Version:         a.version,
		SubscriptionIds: ids,
	}
}
