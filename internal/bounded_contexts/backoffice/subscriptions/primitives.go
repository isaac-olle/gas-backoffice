package subscriptions

type SubscriptionPrimitives struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	OwnerId  string                 `json:"owner_id"`
	IsSystem bool                   `json:"is_system"`
	Members  []MembershipPrimitives `json:"members"`
	Version  int                    `json:"version"`
}

type MembershipPrimitives struct {
	UserId           string `json:"user_id"`
	SubscriptionRole string `json:"subscription_role"`
}

func (s *Subscription) ToPrimitives() SubscriptionPrimitives {
	members := make([]MembershipPrimitives, 0, len(s.members))
	for id, role := range s.members {
		members = append(members, MembershipPrimitives{UserId: id.String(), SubscriptionRole: role.String()})
	}
	return SubscriptionPrimitives{Id: s.id.String(), Name: s.name, OwnerId: s.ownerId.String(), IsSystem: s.isSystem, Members: members, Version: s.version}
}
