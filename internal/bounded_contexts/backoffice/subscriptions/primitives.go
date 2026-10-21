package subscriptions

type SubscriptionPrimitives struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	CreatedBy string   `json:"created_by"`
	MemberIds []string `json:"member_ids"`
	Version   int      `json:"version"`
}

func (s *Subscription) ToPrimitives() SubscriptionPrimitives {
	members := make([]string, 0, s.memberIds.Size())
	for _, id := range s.memberIds.Values() {
		members = append(members, id.String())
	}
	return SubscriptionPrimitives{Id: s.id.String(), Name: s.name, CreatedBy: s.createdBy.String(), MemberIds: members, Version: s.version}
}
